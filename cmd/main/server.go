package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IlianBuh/GraphQL/internal/app"
	"github.com/IlianBuh/GraphQL/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

const defaultPort = "8080"

func main() {
	cfg := config.New()

	log := setUpLogger(cfg.Env)

	application := app.New(log, cfg.Port, cfg.SSOClient)

	application.GraphQL.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop
	log.Info("stopped signal was received", slog.Any("signal", sign))

	application.GraphQL.Stop()
}

func setUpLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return nil
}
