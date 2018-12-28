package helper

import (
	"time"
	"udabayar-go-api-di/config"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWT struct {
	Config *config.Config
}

type Claims struct {
	UserId uint
	Scope  string
	jwt.StandardClaims
}

func (h *JWT) SignJWT(id uint, scope string) (string, error) {
	in6h := time.Now().Add(time.Duration(24) * time.Hour).Unix()

	claims := &Claims{}
	claims.UserId = id
	claims.Scope = "api:" + scope
	claims.ExpiresAt = in6h

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := []byte(h.Config.Key)
	tokenString, err := token.SignedString(key)

	return tokenString, err
}

func (h *JWT) ParseJWT(token string, baseToken *Claims) (*jwt.Token, error) {
	parseJWT, err := jwt.ParseWithClaims(token, baseToken, func(token *jwt.Token) (interface{}, error) {
		key := []byte(h.Config.Key)
		return key, nil
	})

	return parseJWT, err
}

func NewJWTHelper(config *config.Config) *JWT {
	return &JWT{
		Config: config,
	}
}
