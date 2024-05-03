package server

import (
	"kurs-server/domain"
	"kurs-server/infrastructure/server/controllers"
	"kurs-server/infrastructure/tokens"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GinServer struct {
	logger    *zap.Logger
	engine    *gin.Engine
	cases     *domain.UseCases
	tokenizer *tokens.JwtTokenizer
}

func NewGinServer(logger *zap.Logger, uc *domain.UseCases, tokenizer *tokens.JwtTokenizer) *GinServer {
	return &GinServer{
		cases:     uc,
		logger:    logger,
		tokenizer: tokenizer,
		engine:    gin.Default(),
	}
}

func (s *GinServer) SetupRoutes() {
	ctr := controllers.CreateNewController(s.logger, s.cases, s.tokenizer)

	s.engine.GET("/rating/create", ctr.CreateRating) // TODO: implement this

	s.engine.POST("/categories/create", ctr.CreateCategory) // TODO: should create details for them in one request
	s.engine.GET("/categories/get", ctr.GetCategories)
	s.engine.DELETE("/categories/delete", ctr.DeleteCategory)

	s.engine.POST("/products/create", ctr.CreateNewProduct)
	s.engine.DELETE("/products/delete", ctr.DeleteProduct)
	s.engine.GET("/products/get", ctr.GetProducts)
	s.engine.PATCH("/products/edit", ctr.EditProduct)

	s.engine.POST("/users/signup", ctr.CreateNewUser)
	s.engine.POST("/users/signin", ctr.SignIn)
}

func (s *GinServer) Run(addr ...string) {
	err := s.engine.Run(addr...)
	if err != nil {
		s.logger.Error("Unable to run server", zap.Error(err))
		panic(err)
	}
}
