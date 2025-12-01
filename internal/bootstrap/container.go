package bootstrap

import (
	"fmt"

	"github.com/go-chi/chi/v5"
)

type Container struct {
	Status int
	Router *chi.Mux
}

// CreateServerContainer initializes the application container using Wire dependency injection
func CreateServerContainer() (*Container, error) {
	return InitializeContainer()
}

func (c *Container) Close() {
	c.Status = 0
	fmt.Println("Container closed")
}
