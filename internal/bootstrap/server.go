package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/haidang666/go-app/internal/config"
	"github.com/haidang666/go-app/pkg/logger"
)

func StartRestAPI(ctx context.Context, cfg *config.Config, router *chi.Mux) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.L().Infof("listening on :%d", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.L().Info("shutting down server...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		return nil
	case err := <-errCh:
		return err
	}
}

