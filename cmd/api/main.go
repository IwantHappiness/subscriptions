package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	postgresinfra "github.com/IwantHappiness/subscriptions/internal/infrastructure/postgres"
	postgresrepo "github.com/IwantHappiness/subscriptions/internal/repository/postgres"
	transporthttp "github.com/IwantHappiness/subscriptions/internal/transport/http"
	swaggerdocs "github.com/IwantHappiness/subscriptions/internal/transport/http/docs"
	httphandlers "github.com/IwantHappiness/subscriptions/internal/transport/http/handlers"
	"github.com/IwantHappiness/subscriptions/internal/usecase/subscription"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := loadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := postgresinfra.Open(ctx, cfg.DatabaseDSN)
	if err != nil {
		logger.Error("open postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	subRepo := postgresrepo.New(pool)
	subUsecase := subscription.NewService(subRepo)
	subHandler := httphandlers.NewSubHandler(subUsecase, logger)
	docsHandler := swaggerdocs.NewHandler()
	router := transporthttp.NewRouter(subHandler, docsHandler, logger)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown http server", "error", err)
		}
	}()

	logger.Info("http server started", "addr", cfg.HTTPAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("listen and serve", "error", err)
		os.Exit(1)
	}
}

type config struct {
	HTTPAddr    string
	DatabaseDSN string
}

func loadConfig() config {
	cfg := config{
		HTTPAddr:    envOrDefault("HTTP_ADDR", ":8080"),
		DatabaseDSN: envOrDefault("DATABASE_DSN", "postgres://postgres:sub-test-123@localhost:5432/test-sub?sslmode=disable"),
	}

	if cfg.DatabaseDSN == "" {
		panic(fmt.Errorf("DATABASE_DSN is required"))
	}

	return cfg
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
