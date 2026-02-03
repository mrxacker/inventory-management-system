package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/config"
	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/database"
	"github.com/mrxacker/inventory-management-system/services/auth-service/internal/handler"
	"github.com/mrxacker/inventory-management-system/shared/logger"
)

func StartApp(ctx context.Context) error {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
		return err
	}

	log := setupLogger(cfg)

	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database", "error", err)
		return err
	}
	defer db.Close()

	productHandler := handler.NewProductHandler()

	err = StartServer(ctx, cfg, productHandler, log)
	if err != nil {
		log.Fatal("Failed to start server", "error", err)
		return err
	}
	return nil
}

func setupLogger(cfg *config.Config) logger.Logger {
	log := logger.NewLoggerWithConfig(logger.LogConfig{
		Level:       cfg.Log.Level,
		Environment: cfg.Server.Environment,
		OutputPaths: cfg.Log.OutputPaths,
		ErrorPaths:  cfg.Log.ErrorPaths,
	})
	defer log.Sync()

	log.Info("Starting Product Service",
		"version", "1.0.0",
		"environment", cfg.Server.Environment,
		"port", cfg.Server.Port,
	)
	return log
}

func StartServer(ctx context.Context, cfg *config.Config, productHandler *handler.ProductHandler, log logger.Logger) error {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := handler.SetupRouter(productHandler, cfg, log)

	// Create HTTP server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(cfg.Server.IdleTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	errChan := make(chan error, 1)

	go func() {
		log.Info("Starting server", "port", cfg.Server.Port, "environment", cfg.Server.Environment)
		if err := srv.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}

	}()

	select {
	case <-ctx.Done():
		log.Info("Shutdown signal received")
	case err := <-errChan:
		log.Info("Server error", "error", err)
		return err
	}

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited successfully")

	return nil
}
