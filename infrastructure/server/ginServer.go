package server

import (
	"kurs-server/domain"
	"kurs-server/infrastructure/server/controllers"
	"kurs-server/infrastructure/server/middleware"
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

	guard := middleware.Guard(s.cases, s.tokenizer, s.logger)
	// create orders, order groups

	s.engine.POST("/reviews/create", ctr.CreateReview) //TODO: implement this (for authorized user)
	// s.engine.DELETE("/reviews/delete", ctr.DeleteReview) // POHUI??

	s.engine.POST("/rating/create", ctr.CreateRating)  // TODO: implement this (for authorized user)
	s.engine.DELETE("rating/delete", ctr.DeleteRating) //TODO: implement this (for authorized user)

	s.engine.POST("/categories/create", ctr.CreateCategory)   // works fine (for admin)
	s.engine.GET("/categories/get", ctr.GetCategories)        // works fine (for everyone)
	s.engine.DELETE("/categories/delete", ctr.DeleteCategory) // works fine (for admin)

	s.engine.POST("/products/create", ctr.CreateNewProduct)  //somehow works (for admin)
	s.engine.DELETE("/products/delete", ctr.DeleteProduct)   //works fine
	s.engine.GET("/products/get", guard, ctr.GetProducts)    //works fine
	s.engine.PATCH("/products/edit", ctr.EditProduct)        //works fine
	s.engine.GET("/products/info", guard, ctr.GetProuctInfo) //works fine

	s.engine.POST("/users/signup", ctr.CreateNewUser) //works fine i guess (for everyone)
	s.engine.POST("/users/signin", ctr.SignIn)        //works fine i guess (for everyone)
	s.engine.POST("/users/validate", ctr.ValidateUser)
}

func (s *GinServer) Run(addr ...string) {
	err := s.engine.Run(addr...)
	if err != nil {
		s.logger.Error("Unable to run server", zap.Error(err))
		panic(err)
	}
}
