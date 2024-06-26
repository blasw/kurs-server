package middleware

import (
	"errors"
	"fmt"
	"kurs-server/domain"
	tokens2 "kurs-server/infrastructure/tokens"
	"kurs-server/structs"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func Guard(uc *domain.UseCases, jwtTokenizer *tokens2.JwtTokenizer, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		access_token, err := c.Cookie("access_token")
		if err != nil {
			logger.Debug("Unable to find access_token")
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		refresh_token, err := c.Cookie("refresh_token")
		if err != nil {
			logger.Debug("unable to find refresh_token")
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		res, err := ValidateUser(uc, jwtTokenizer, access_token, refresh_token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("ID", res.UserID)
		c.Set("Username", res.Username)
		c.Set("Role", res.Role)
		c.Next()
	}
}

type validationResults struct {
	Access_Token  string
	Refresh_Token string
	UserID        uint
	Username      string
	Role          string
}

func ValidateUser(uc *domain.UseCases, tokenizer *tokens2.JwtTokenizer, access_token string, refresh_token string) (*validationResults, error) {
	accessClaims, err := tokenizer.ParseAccessToken(access_token)
	if err == nil && accessClaims != nil {
		//access token is valid
		user_id, err := strconv.Atoi(accessClaims.ID)
		if err != nil {
			return nil, err
		}

		res := &validationResults{
			Access_Token:  "",
			Refresh_Token: "",
			UserID:        uint(user_id),
			Username:      accessClaims.Username,
			Role:          accessClaims.Role,
		}

		return res, nil
	}

	// access token is invalid
	_, err = tokenizer.ParseRefreshToken(refresh_token)
	if err != nil {
		//refresh token is invalid
		return nil, errors.New("invalid tokens")
	}

	//access token is invalid, refresh token is valid
	user, err := uc.Users().GetUserByRefreshToken(refresh_token)
	if err != nil {
		return nil, errors.New("invalid tokens")
	}

	newRefreshToken, err := tokenizer.NewRefreshToken(jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	if err != nil {
		return nil, errors.New("")
	}

	uc.Users().UpdateUserRefreshToken(user.Username, newRefreshToken)

	newAccessToken, err := tokenizer.NewAccessToken(structs.JwtUserClaims{
		ID:       fmt.Sprint(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	res := &validationResults{
		UserID:        user.ID,
		Username:      user.Username,
		Role:          user.Role,
		Access_Token:  newAccessToken,
		Refresh_Token: newRefreshToken,
	}
	return res, nil
}
