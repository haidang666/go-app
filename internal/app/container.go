package app

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/haidang666/go-app/internal/bootstrap"
	"github.com/haidang666/go-app/internal/config"
	infrastructure "github.com/haidang666/go-app/internal/infrastructure/repository"
)

type Container struct {
	Status int
	Router *chi.Mux
}

func CreateServerContainer(ctx context.Context, cfg *config.Config) (*Container, error) {

	r := bootstrap.BootstrapHandler(
		bootstrap.HandlerBootstrapArgs{
			Repositories: bootstrap.Repositories{
				UserRepository: infrastructure.NewUserRepository(),
			},
		},
	)

	return &Container{
		Status: 1,
		Router: r,
	}, nil
}

func (c *Container) Close() {
	c.Status = 0
	fmt.Println("Container closed")
}
