package session

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/request"

	infoblog "gitlab.com/InfoBlogFriends/server"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

const (
	privKeyPath = "./keys/private"
	pubKeyPath  = "./keys/public"
)

type JWTKeys struct {
	signKey     *rsa.PrivateKey
	verifyKey   *rsa.PublicKey
	signBytes   []byte
	verifyBytes []byte
	logger      *zap.Logger
}

func NewJWTKeys(logger *zap.Logger) (*JWTKeys, error) {
	j := &JWTKeys{
		logger: logger,
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

func (j *JWTKeys) CreateToken(u infoblog.User) (string, error) {
	signKey, err := j.GetSignKey()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": time.Now().UTC().Add(time.Minute * 20).Unix(),
		"id":  u.ID,
	})

	return token.SignedString(signKey)
}

func (j *JWTKeys) ExtractToken(r *http.Request) (*infoblog.User, error) {
	verifyKey, err := j.GetVerifyKey()
	if err != nil {
		return nil, err
	}
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
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
	if v, ok := c["id"]; ok {
		uid = int64(v.(float64))
	}
	u := &infoblog.User{
		ID: uid,
	}

	return u, nil

}
