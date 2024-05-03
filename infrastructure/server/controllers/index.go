package controllers

import (
	"go.uber.org/zap"
	"kurs-server/domain"
	"kurs-server/infrastructure/tokens"
)

type Controller struct {
	logger    *zap.Logger
	cases     *domain.UseCases
	tokenizer *tokens.JwtTokenizer
}

func CreateNewController(logger *zap.Logger, uc *domain.UseCases, tokenizer *tokens.JwtTokenizer) *Controller {
	return &Controller{
		logger,
		uc,
		tokenizer,
	}
}
