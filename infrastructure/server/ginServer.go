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

	s.engine.PATCH("/orderGroups/mark", ctr.MarkOrderGroup)
	s.engine.DELETE("/orderGroups/delete", ctr.DeleteOrderGroup)

	s.engine.GET("/orders/getAllGroups", ctr.GetAllOrderGroups)
	s.engine.GET("/orders/getGroups", guard, ctr.GetOrderGroups)   //works fine
	s.engine.POST("/orders/confirm", guard, ctr.ConfirmOrderGroup) //works fine
	s.engine.GET("/orders/getCart", guard, ctr.GetCart)            //works fine
	s.engine.POST("/orders/create", guard, ctr.CreateOrder)        //works fine
	s.engine.DELETE("/orders/delete", guard, ctr.DeleteOrder)      //works fine

	s.engine.GET("/details/get", ctr.GetDetailsValues) // works fine

	s.engine.POST("/reviews/create", guard, ctr.CreateReview) //works fine
	s.engine.GET("/reviews/get", ctr.GetReviews)              //works fine

	s.engine.POST("/rating/create", ctr.CreateRating) // works fine

	s.engine.POST("/categories/create", ctr.CreateCategory)   // works fine
	s.engine.GET("/categories/get", ctr.GetCategories)        // works fine
	s.engine.DELETE("/categories/delete", ctr.DeleteCategory) // works fine

	s.engine.POST("/products/create", ctr.CreateNewProduct)           //works fine
	s.engine.DELETE("/products/delete", ctr.DeleteProduct)            //works fine
	s.engine.GET("/products/get", guard, ctr.GetProducts)             //works fine BUT without details filter
	s.engine.POST("/products/search", guard, ctr.FilteredGetProducts) //works fine WITH details filter
	s.engine.PATCH("/products/edit", ctr.EditProduct)                 //works fine
	s.engine.GET("/products/info", guard, ctr.GetProductInfo)         //works fine

	s.engine.POST("/users/signup", ctr.CreateNewUser)  //works fine
	s.engine.POST("/users/signin", ctr.SignIn)         //works fine
	s.engine.POST("/users/validate", ctr.ValidateUser) //works fine
}

func (s *GinServer) Run(addr ...string) {
	err := s.engine.Run(addr...)
	if err != nil {
		s.logger.Error("Unable to run server", zap.Error(err))
		panic(err)
	}
}
