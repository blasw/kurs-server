package application

import (
	"kurs-server/application/utility"
	"kurs-server/domain"
	"kurs-server/infrastructure/db"
	"kurs-server/infrastructure/server"
	"kurs-server/infrastructure/tokens"
	"os"

	"go.uber.org/zap"
)

func init() {
	utility.CheckEnv()
}

func Run(addr ...string) {
	logger, _ := zap.NewDevelopment()
	logger.Info("Starting kurs-server...")

	logger.Debug("Connecting to database...")
	store := db.NewPostgreStore(logger, os.Getenv("DB_ADDR"))
	logger.Debug("Database connected...")

	uc := domain.NewUseCases(store.DB())

	// Filling db with some test data
	utility.AddTestData(uc)

	tokenizer := tokens.CreateJwtTokenizer(logger, os.Getenv("SECRET"))

	logger.Debug("Creating servre instance...")
	serverInstance := server.NewGinServer(logger, uc, tokenizer)
	serverInstance.SetupRoutes()

	logger.Debug("Running server instance...")
	serverInstance.Run(addr...)
}
