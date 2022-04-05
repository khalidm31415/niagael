package main

import (
	"cart/entity"
	"cart/graphql"
	"cart/graphql/schema"
	"cart/service"
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	graphql_go "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
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

	cartService := service.NewCartService(db)
	productServiceURL := "http://product_api:2719/product/graphql"
	// productServiceURL := "http://localhost:2719/product/graphql"
	productService := service.NewProductService(productServiceURL)
	resolver := graphql.NewResolver(cartService, productService)

	schemaString, err := schema.String()
	if err != nil {
		zap.S().Fatal(err)
	}

	graphqlSchema := graphql_go.MustParseSchema(schemaString, resolver)

	app := fiber.New()
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/cart/graphql", graphqlHandler(&relay.Handler{Schema: graphqlSchema}))

	app.Get("/cart/graphiql", func(c *fiber.Ctx) error {
		return c.SendFile("web/graphiql.html")
	})

	port := env["PORT"]
	app.Listen(fmt.Sprintf(":%s", port))
}

func setupLogger() {
	configLogger := zap.NewDevelopmentConfig()
	configLogger.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	configLogger.DisableStacktrace = true
	logger, _ := configLogger.Build()
	zap.ReplaceGlobals(logger)
}

func graphqlHandler(h http.Handler) fiber.Handler {
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-user-id")
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
