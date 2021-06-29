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
		"ExpiresAt": time.Now().UTC().Add(time.Minute * 20).Unix(),
		"ID":        u.ID,
		"Email":     u.Email,
		"Phone":     u.Phone,
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

	if token.Valid {
		c := token.Claims.(jwt.MapClaims)
		u := &infoblog.User{
			ID:    c["ID"].(int64),
			Phone: c["Phone"].(string),
			Email: c["Email"].(string),
		}

		return u, nil
	}

	return nil, errors.New("invalid token")

}
