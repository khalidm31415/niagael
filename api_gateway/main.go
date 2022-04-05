package main

import (
	"api_gateway/entity"
	"api_gateway/router"
	"api_gateway/service"
	"fmt"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	setupLogger()
	zap.S().Info("Bismillah")
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		zap.S().Fatal(err)
	}
	zap.S().Debug(zap.Any("env", env))

	dsn := env["POSTGRES_DSN"]
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		zap.S().Fatal(err)
	}
	if err := entity.Migrate(db); err != nil {
		zap.S().Fatal(err)
	}

	authService := service.NewAuthService(db)
	r := router.NewRouter(authService)
	port := env["PORT"]
	r.Run(fmt.Sprintf(":%s", port))
}

func setupLogger() {
	configLogger := zap.NewDevelopmentConfig()
	configLogger.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	configLogger.DisableStacktrace = true
	logger, _ := configLogger.Build()
	zap.ReplaceGlobals(logger)
}
