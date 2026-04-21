package server

import (
	"backend/config"
	"backend/internal/adapters/bcrypt"
	api "backend/internal/adapters/http"
	"backend/internal/adapters/jwt"
	"backend/internal/adapters/postgres"
	"backend/internal/domain/auth"
	"backend/internal/domain/user"
	"backend/pkg/database"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port     string
	db       postgres.DBService
	handlers *api.Handlers
}

func NewServer() *http.Server {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := database.RunMigrations(cfg.DBConnString, cfg.MigrationsPath, logger); err != nil {
		logger.Error("Migrations failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db := postgres.New(cfg.DBConnString, cfg.DBName)

	// Repositories
	authRepository := postgres.NewAuthRepository(db.Pool())
	userRepository := postgres.NewUserRepository(db.Pool())

	tokenProvider, err := jwt.NewTokenProvider(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	if err != nil {
		log.Fatalf("Failed to initialize token provider: %v", err)
	}
	hasher := bcrypt.NewPasswordHasher()

	// Services
	authService := auth.NewService(authRepository, tokenProvider, hasher)
	userService := user.NewService(userRepository)

	// Handlers
	handlers := &api.Handlers{
		Auth: api.NewAuthHandler(authService, logger),
		User: api.NewUserHandler(userService, logger),
	}

	app := &Server{
		port:     cfg.Port,
		db:       db,
		handlers: handlers,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", app.port),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
