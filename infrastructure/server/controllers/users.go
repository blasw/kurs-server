package controllers

import (
	"fmt"
	"kurs-server/domain/entities"
	"kurs-server/infrastructure/server/middleware"
	"kurs-server/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type createUserDto struct {
	Username string `json:"username" binding:"required,min=4,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
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

	//creating new user
	newUser := &entities.User{
		Username:     dto.Username,
		Password:     string(hash),
		Email:        dto.Email,
		Role:         "User",
		RefreshToken: refreshToken,
	}

	newUserID, err := ctr.cases.Users().Create(newUser)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

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

	//setting cookies
	c.SetCookie("access_token", accessToken, 3600*24, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", false, true)

	resp := gin.H{"username": newUser.Username, "userid": newUserID, "role": "User", "access_token": accessToken, "refresh_token": refreshToken}

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

	storedUser, err := ctr.cases.Users().GetByUsername(dto.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
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

	c.SetCookie("access_token", accessToken, 3600*24, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", false, true)

	ctr.cases.Users().UpdateUserRefreshToken(dto.Username, refreshToken)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken, "userID": storedUser.ID, "username": storedUser.Username, "role": storedUser.Role})
}

type ValidateUserDto struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (ctr *Controller) ValidateUser(c *gin.Context) {
	var dto ValidateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rToken := dto.RefreshToken
	aToken := dto.AccessToken

	res, err := middleware.ValidateUser(ctr.cases, ctr.tokenizer, aToken, rToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
