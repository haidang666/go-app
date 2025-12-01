//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/haidang666/go-app/internal/domain/contract"
	authUseCase "github.com/haidang666/go-app/internal/domain/use_case/auth"
	"github.com/haidang666/go-app/internal/infrastructure/http/handlers/auth"
	"github.com/haidang666/go-app/internal/infrastructure/http/router"
	infrastructure "github.com/haidang666/go-app/internal/infrastructure/repository"
)

// Providers for the application container
var ProviderSet = wire.NewSet(
	ProvideUserRepository,
	ProvideSignUpUseCase,
	ProvideAuthHandler,
	ProvideRouter,
	ProvideContainer,
)

// ProvideUserRepository provides the user repository implementation
func ProvideUserRepository() contract.UserRepository {
	return infrastructure.NewUserRepository()
}

// ProvideSignUpUseCase provides the sign up use case
func ProvideSignUpUseCase(userRepo contract.UserRepository) *authUseCase.SignUpUseCase {
	return authUseCase.NewSignUpUseCase(userRepo)
}

// ProvideAuthHandler provides the auth handler
func ProvideAuthHandler(signUpUseCase *authUseCase.SignUpUseCase) *auth.AuthHandler {
	return auth.NewAuthHandler(auth.NewAuthHandlerArgs{
		SignUpUseCase: signUpUseCase,
	})
}

// ProvideRouter provides the chi router with all routes registered
func ProvideRouter(authHandler *auth.AuthHandler) *chi.Mux {
	return router.NewRouter(router.NewRouterArgs{
		AuthHandler: authHandler,
	})
}

// ProvideContainer provides the application container
func ProvideContainer(r *chi.Mux) *Container {
	return &Container{
		Status: 1,
		Router: r,
	}
}

// InitializeContainer initializes and returns the application container
// This function is implemented by the wire code generator
func InitializeContainer() (*Container, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
