package main

import (
	"context"
	"fmt"
	"net/http"
	"product/entity"
	"product/graphql"
	"product/graphql/schema"
	"product/service"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	graphql_go "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
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

	dsn := env["MYSQL_DSN"]
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		zap.S().Fatal(err)
	}
	if err := entity.Migrate(db); err != nil {
		zap.S().Fatal(err)
	}

	productService := service.NewProductService(db)
	resolver := graphql.NewResolver(productService)

	schemaString, err := schema.String()
	if err != nil {
		zap.S().Fatal(err)
	}

	graphqlSchema := graphql_go.MustParseSchema(schemaString, resolver)

	app := fiber.New()
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/product/graphql", graphqlHandler(&relay.Handler{Schema: graphqlSchema}))

	app.Get("/product/graphiql", func(c *fiber.Ctx) error {
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
