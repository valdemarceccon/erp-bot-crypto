package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

func notImplemented(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func main() {
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")

	dbConfig := repository.PostgresConfigFromEnv()

	db, err := sql.Open("pgx", dbConfig.ToString())

	if err != nil {
		log.Fatal(err)
	}

	if port == "" {
		port = "8000"
	}

	if jwtSecret == "" {
		jwtSecret = "some-secret"
	}

	userRepo := repository.NewUserRepositoryPsql(db)

	authControler := controller.NewJwtAuthController(userRepo, controller.WithHS256Secret(jwtSecret))
	userController := controller.NewUserController(userRepo)

	app := fiber.New()

	authGroup := app.Group("/auth")
	authGroup.Post("/login", authControler.LoginHandler)
	authGroup.Get("/logout", notImplemented)
	authGroup.Get("/refresh", notImplemented)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
	})

	userGroup := app.Group("/user")
	userGroup.Use(jwtMiddleware)

	userGroup.Get("/", userController.ListUsers)

	app.Listen(":" + port)
}

type Resp struct {
	Message string `json:"message"`
}
