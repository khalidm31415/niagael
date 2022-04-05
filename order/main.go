package main

import (
	"context"
	"fmt"
	"net/http"
	"order/entity"
	"order/graphql"
	"order/graphql/schema"
	"order/service"

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

	cartServiceURL := "http://cart_api:2720/cart/graphql"
	// cartServiceURL := "http://localhost:2720/cart/graphql"
	cartService := service.NewCartService(cartServiceURL)
	orderService := service.NewOrderService(db, cartService)
	resolver := graphql.NewResolver(orderService)

	schemaString, err := schema.String()
	if err != nil {
		zap.S().Fatal(err)
	}

	graphqlSchema := graphql_go.MustParseSchema(schemaString, resolver)

	app := fiber.New()
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/order/graphql", graphqlHandler(&relay.Handler{Schema: graphqlSchema}))

	app.Get("/order/graphiql", func(c *fiber.Ctx) error {
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
