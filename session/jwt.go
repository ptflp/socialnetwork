package session

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/cache"

	"gitlab.com/InfoBlogFriends/server/hasher"

	"github.com/dgrijalva/jwt-go/request"

	infoblog "gitlab.com/InfoBlogFriends/server"
	req "gitlab.com/InfoBlogFriends/server/request"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

const (
	privKeyPath = "./keys/private"
	pubKeyPath  = "./keys/public"

	Month = 30 * Day
	Day   = 24 * time.Hour

	RefreshTokenKey = "refresh_token"
)

type JWTKeys struct {
	*decoder.Decoder
	signKey     *rsa.PrivateKey
	verifyKey   *rsa.PublicKey
	signBytes   []byte
	verifyBytes []byte
	logger      *zap.Logger
	cache       cache.Cache
}

func NewJWTKeys(logger *zap.Logger, cache cache.Cache) (*JWTKeys, error) {
	j := &JWTKeys{
		Decoder: decoder.NewDecoder(),
		logger:  logger,
		cache:   cache,
	}
	err := j.ReadKeys()
	if err != nil {
		return nil, err
	}

	return j, nil
}

func (j *JWTKeys) ReadKeys() error {
	var err error

	if j.signBytes == nil {
		j.signBytes, err = ioutil.ReadFile(privKeyPath)
		if err != nil {
			j.logger.Error("read private key err", zap.Error(err))
			return err
		}
	}
	j.signKey, err = jwt.ParseRSAPrivateKeyFromPEM(j.signBytes)
	if err != nil {
		j.logger.Error("read private key err", zap.Error(err))
		return err
	}

	if j.verifyBytes == nil {
		j.verifyBytes, err = ioutil.ReadFile(pubKeyPath)
		if err != nil {
			j.logger.Error("read private key err", zap.Error(err))
			return err
		}
	}

	j.verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(j.verifyBytes)
	if err != nil {
		j.logger.Error("read private key err", zap.Error(err))
		return err
	}

	return nil
}

func (j *JWTKeys) GetSignKey() (*rsa.PrivateKey, error) {
	if j.signKey == nil {
		j.logger.Error("retrieve signKey err")
		return nil, errors.New("retrieve signKey err")
	}

	return j.signKey, nil
}

func (j *JWTKeys) GetVerifyKey() (*rsa.PublicKey, error) {
	if j.verifyKey == nil {
		j.logger.Error("retrieve verifyKey err")
		return nil, errors.New("retrieve verifyKey err")
	}

	return j.verifyKey, nil
}

type JWTAuth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	Token string `json:"token"`
	UID   int64  `json:"uid"`
	UUID  string `json:"uuid"`
}

func (j *JWTKeys) GenerateAuthTokens(u *infoblog.User) (*req.AuthTokenData, error) {
	if len(u.UUID) < 40 {
		return nil, errors.New("wrong user")
	}
	access, err := j.CreateAccessToken(*u)
	if err != nil {
		return nil, err
	}
	refresh, err := j.CreateRefreshToken(access, u)
	if err != nil {
		return nil, err
	}

	authToken := req.AuthTokenData{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	err = j.MapStructs(&authToken.User, u)
	if err != nil {
		return nil, err
	}

	return &authToken, err
}

func (j *JWTKeys) CreateAccessToken(u infoblog.User) (string, error) {
	token, err := j.GenerateToken(jwt.MapClaims{
		"exp":  time.Now().UTC().Add(time.Hour * 50).Unix(),
		"uid":  u.ID,
		"uuid": u.UUID,
	})

	return token, err
}

func (j *JWTKeys) CreateRefreshToken(accessToken string, u *infoblog.User) (string, error) {
	refreshToken := hasher.NewSHA256([]byte(accessToken))

	key := strings.Join([]string{RefreshTokenKey, strconv.Itoa(int(u.ID)), refreshToken}, ":")
	j.cache.Set(key, u, 2*Month)

	return j.GenerateToken(jwt.MapClaims{
		"refresh_token": refreshToken,
		"exp":           time.Now().UTC().Add(2 * Month).Unix(),
		"uid":           u.ID,
		"uuid":          u.UUID,
	})
}

func (j *JWTKeys) GenerateToken(m jwt.MapClaims) (string, error) {
	signKey, err := j.GetSignKey()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, m)
	if err != nil {
		return "", err
	}

	return token.SignedString(signKey)
}

func (j *JWTKeys) ExtractRefreshToken(rawToken string) (*RefreshToken, error) {
	verifyKey, err := j.GetVerifyKey()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	c := token.Claims.(jwt.MapClaims)

	v, ok := c["exp"]
	if !ok {
		return nil, errors.New("jwt map claims err: exp")
	}

	exp := int64(v.(float64))
	now := time.Now().Unix()
	if now > exp {
		return nil, errors.New("refresh token expired")
	}

	refreshToken, ok := c["refresh_token"]
	if !ok {
		return nil, errors.New("jwt map claims err: refresh_token")
	}

	uid, ok := c["uid"]
	if !ok {
		return nil, errors.New("jwt map claims err: uid")
	}

	uuid, ok := c["uuid"]
	if !ok {
		return nil, errors.New("jwt map claims err: uuid")
	}

	return &RefreshToken{
		Token: refreshToken.(string),
		UID:   int64(uid.(float64)),
		UUID:  uuid.(string),
	}, nil
}

func (j *JWTKeys) ExtractAccessToken(r *http.Request) (*infoblog.User, error) {
	verifyKey, err := j.GetVerifyKey()
	if err != nil {
		return nil, err
	}
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	c := token.Claims.(jwt.MapClaims)
	var uid, exp int64

	if v, ok := c["ExpiresAt"]; ok {
		exp = int64(v.(float64))
	}
	if v, ok := c["exp"]; ok {
		exp = int64(v.(float64))
	}

	now := time.Now().Unix()
	if now > exp {
		return nil, errors.New("token expired")
	}

	if v, ok := c["ID"]; ok {
		uid = int64(v.(float64))
	}
	var uuid string
	if v, ok := c["uuid"]; ok {
		uuid = v.(string)
	}

	if v, ok := c["uid"]; ok {
		uid = int64(v.(float64))
	}
	u := &infoblog.User{
		ID:   uid,
		UUID: uuid,
	}

	return u, nil

}
