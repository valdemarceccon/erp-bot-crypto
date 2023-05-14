package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/shopspring/decimal"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/middleware"
	"github.com/valdemarceccon/crypto-bot-erp/backend/migrations"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"
)

func notImplemented(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func runMigrations(db *sql.DB) {
	goose.SetBaseFS(migrations.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	shouldMigrate := os.Getenv("ENABLE_MIGRATIONS")

	dbConfig := store.PostgresConfigFromEnv()

	db, err := sql.Open("pgx", dbConfig.String())

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if strings.ToLower(shouldMigrate) == "true" {
		runMigrations(db)
	}

	if port == "" {
		port = "8000"
	}

	if jwtSecret == "" {
		jwtSecret = "some-secret"
	}

	config := &model.AppConfig{
		Commission: decimal.RequireFromString("0.4"),
	}

	userStore := store.NewUserPsql(db)
	roleStore := store.NewRolePsql(db)
	apiStore := store.NewApiKeyPsql(db)

	authControler := controller.NewJwtAuthController(userStore, controller.WithHS256Secret(jwtSecret))
	userController := controller.NewUserController(userStore, roleStore, apiStore, config)
	dataCollectorController := controller.NewDataCollector(userStore, apiStore)

	authMiddleware := middleware.NewAuthMiddleware(userStore, roleStore, jwtSecret)

	guard := controller.NewGuards(roleStore)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Println(err)
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = ctx.Status(code).JSON(schema.MessageResponse{
				Message: e.Message,
			})

			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(schema.MessageResponse{Message: "Internal Server Error"})
			}

			// Return from handler
			return nil
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	authGroup := app.Group("/auth")
	authGroup.Post("/login", authControler.LoginHandler)
	authGroup.Post("/register", authControler.RegisterHandler)
	authGroup.Get("/logout", notImplemented)
	authGroup.Get("/refresh", notImplemented)

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
	})

	userGroup := app.Group("/user")
	userGroup.Use(jwtMiddleware)
	userGroup.Use(authMiddleware.UserExists)

	userGroup.Get("/", guard.WithPermission(model.ListUsersPermission, userController.ListUsers))
	userGroup.Get("/api_keys", userController.ListApiKeys)
	userGroup.Get("/api_keys/all", guard.WithPermission(model.ListApiKeysPermission, userController.ListAllApiKeys))
	userGroup.Post("/api_keys", userController.AddApiKey)
	userGroup.Patch("/api_keys/client-toggle/:apiKeyId", userController.ClientToggleApiKey)
	userGroup.Patch("/api_keys/admin-toggle/:userId/:apiKeyId", guard.WithPermission(model.WriteApiKeysPermission, userController.AdminToggleApiKey))
	userGroup.Get("/comission/:userId", guard.WithPermission(model.GetUserCommission, userController.CalculateComission))
	userGroup.Get("/me", userController.Me)

	collectorGroup := app.Group("/collector")
	collectorGroup.Use(jwtMiddleware)
	collectorGroup.Use(authMiddleware.UserExists)
	//{{host}}/collector/2023-05-06/2023-05-06
	collectorGroup.Post("/:startDate/:endDate/:username?", guard.WithPermission(model.RunDataCollectorPermission, dataCollectorController.RunNow))

	app.Listen(":" + port)
}

type Resp struct {
	Message string `json:"message"`
}
