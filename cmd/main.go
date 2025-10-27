package main

import (
	"context"
	_ "marketplace/docs" // Важно: добавить импорт docs
	"marketplace/internal/configs"
	"marketplace/internal/controller"
	"marketplace/internal/db"
	"marketplace/internal/repository"
	"marketplace/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// @title           Marketplace API
// @version         1.0
// @description     API Сервер для маркетплейса
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  saidnurnasreddinzoda@gmail.com

// @host      localhost:7577
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	log.Info().Msg("Starting up - Start")
	if err := configs.ReadSettings(); err != nil {
		log.Error().Err(err).Msg("Error reading settings" + err.Error())
		return
	}

	log.Info().Any("configs", configs.AppSettings).Msg("Configs")

	dbHost := configs.AppSettings.PostgresParams.Host
	redisAddr := configs.AppSettings.RedisParams.Addr

	if os.Getenv("DOCKER_ENV") == "true" {
		dbHost = "db"
		redisAddr = "redis:6379"
	}

	configs.AppSettings.PostgresParams.Host = dbHost
	configs.AppSettings.RedisParams.Addr = redisAddr

	dbConn, err := db.InitConnection()
	if err != nil {
		log.Error().Err(err).Msg("Error during database connection initialization: " + err.Error())
		return
	}

	repo := repository.NewRepository(dbConn)
	svc := service.NewService(repo, nil)
	ctrl := controller.NewController(svc)

	router := ctrl.InitRoutes()
	srv := &http.Server{
		Addr:    ":" + configs.AppSettings.AppParams.PortRun,
		Handler: router,
	}

	go func() {
		log.Info().Msg("Starting server on port " + configs.AppSettings.AppParams.PortRun)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	if err = db.CloseConnection(dbConn); err != nil {
		log.Error().Err(err).Msg("Error during database connection close")
	} else {
		log.Info().Msg("Database connection closed")
	}

	log.Info().Msg("Server exiting")
}
