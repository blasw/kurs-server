package tokens

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"kurs-server/structs"
)

type JwtTokenizer struct {
	Logger *zap.Logger
	Secret string
}

func CreateJwtTokenizer(logger *zap.Logger, secret string) *JwtTokenizer {
	return &JwtTokenizer{Logger: logger, Secret: secret}
}

func (j *JwtTokenizer) NewAccessToken(userClaims structs.JwtUserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	return accessToken.SignedString([]byte(j.Secret))
}

func (j *JwtTokenizer) NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(j.Secret))
}

func (j *JwtTokenizer) ParseAccessToken(accessToken string) (*structs.JwtUserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &structs.JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedAccessToken.Valid {
		return nil, errors.New("invalid access token provided")
	}

	return parsedAccessToken.Claims.(*structs.JwtUserClaims), nil
}

func (j *JwtTokenizer) ParseRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedRefreshToken.Valid {
		return nil, errors.New("invalid refresh token provided")
	}

	return parsedRefreshToken.Claims.(*jwt.StandardClaims), nil
}
