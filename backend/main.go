package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller"
	"github.com/valdemarceccon/crypto-bot-erp/backend/middleware"
	"github.com/valdemarceccon/crypto-bot-erp/backend/migrations"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
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

	dbConfig := repository.PostgresConfigFromEnv()

	db, err := sql.Open("pgx", dbConfig.ToString())

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

	userRepo := repository.NewUserPsql(db)
	roleRepo := repository.NewRolePsql(db)

	authControler := controller.NewJwtAuthController(userRepo, controller.WithHS256Secret(jwtSecret))
	userController := controller.NewUserController(userRepo, roleRepo)

	authMiddleware := middleware.NewAuthMiddleware(userRepo, roleRepo, jwtSecret)

	app := fiber.New()

	authGroup := app.Group("/auth")
	authGroup.Post("/login", authControler.LoginHandler)
	authGroup.Post("/register", authControler.RegisterHandler)
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
	userGroup.Use(authMiddleware.UserExists)

	userGroup.Get("/", controller.WithPermission(roleRepo, model.ListUsersPermission, userController.ListUsers))
	userGroup.Get("/api_keys", userController.ListApiKeys)
	userGroup.Post("/api_keys", userController.AddApiKey)
	userGroup.Get("/me", userController.Me)

	app.Listen(":" + port)
}

type Resp struct {
	Message string `json:"message"`
}
