package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"kurs-server/domain/entities"
	"kurs-server/structs"
	"net/http"
	"time"
)

type createUserDto struct {
	Username string `json:"username" binding:"required,min=4,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

func (ctr *Controller) CreateNewUser(c *gin.Context) {
	// validating password
	var dto createUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// hashing password
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		ctr.logger.Error("Unable to hash provided password", zap.String("Error: ", err.Error()))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// creating refresh token
	refreshToken, err := ctr.tokenizer.NewRefreshToken(jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	if err != nil {
		ctr.logger.Error("Unable to create new refresh token", zap.String("Error: ", err.Error()))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctr.logger.Debug("refresh token created")

	//creating new user
	newUser := &entities.User{
		Username:     dto.Username,
		Password:     string(hash),
		Email:        dto.Email,
		Role:         "Admin",
		RefreshToken: refreshToken,
	}

	ctr.logger.Debug("new user struct created: ", zap.Any("struct: ", newUser))

	newUserID, err := ctr.cases.Users().Create(newUser)
	if err != nil {
		ctr.logger.Error("Unable to create a new user", zap.String("Error: ", err.Error()))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctr.logger.Debug("new user in db created")

	//creating access token
	accessToken, err := ctr.tokenizer.NewAccessToken(structs.JwtUserClaims{
		ID:       fmt.Sprint(newUserID),
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	if err != nil {
		ctr.logger.Error("Unable to create an access token", zap.String("Error: ", err.Error()))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctr.logger.Debug("new access token created")

	//setting cookies
	c.SetCookie("access_token", accessToken, 3600*24, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", false, true)

	ctr.logger.Debug("cookies set")

	resp := gin.H{"username": newUser.Username, "email": newUser.Email}

	c.JSON(201, resp)
}

type signInDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ctr *Controller) SignIn(c *gin.Context) {
	var dto signInDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctr.logger.Debug("dto: ", zap.Any("struct: ", dto))

	storedUser, err := ctr.cases.Users().GetByUsername(dto.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(dto.Password)); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// at this point username and password are correct
	userClaims := &structs.JwtUserClaims{
		ID:       fmt.Sprint(storedUser.ID),
		Username: storedUser.Username,
		Email:    storedUser.Email,
		Role:     storedUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	accessToken, err := ctr.tokenizer.NewAccessToken(*userClaims)
	if err != nil {
		ctr.logger.Debug("Unable to create new access_token", zap.String("Error: ", err.Error()))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := ctr.tokenizer.NewRefreshToken(jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	if err != nil {
		ctr.logger.Debug("Unable to create new refresh_token", zap.String("Error: ", err.Error()))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctr.logger.Debug("tokens: ", zap.String("tokens: ", accessToken+" <-> "+refreshToken))

	c.SetCookie("access_token", accessToken, 3600*24, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"Message": "Success"})
}

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "Success"})
	}
}
