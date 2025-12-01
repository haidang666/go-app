package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/haidang666/go-app/internal/bootstrap"
	"github.com/haidang666/go-app/internal/config"
	"github.com/haidang666/go-app/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.L().Fatalf("config error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c, err := bootstrap.CreateServerContainer()
	if err != nil {
		logger.L().Fatalf("fail to create server container: %v", err)
	}
	defer c.Close()

	if err := bootstrap.StartRestAPI(ctx, cfg, c.Router); err != nil {
		logger.L().Fatalf("starting server: %v", err)
	}
}
