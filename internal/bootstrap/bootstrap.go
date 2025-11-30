package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/haidang666/go-app/internal/domain/repository"
	authUseCase "github.com/haidang666/go-app/internal/domain/use_case/auth"
	"github.com/haidang666/go-app/internal/infrastructure/http/module_route/auth"
	"github.com/haidang666/go-app/internal/infrastructure/http/router"
)

type HandlerBootstrapArgs struct {
	Repositories Repositories
}

type Repositories struct {
	UserRepository        repository.UserRepository
}

func BootstrapHandler(args HandlerBootstrapArgs) *chi.Mux {
	authHandler := auth.NewAuthHandler(auth.NewAuthHandlerArgs{
		SignUpUseCase: authUseCase.NewSignUpUseCase(args.Repositories.UserRepository),
	})

	r := router.NewRouter(router.NewRouterArgs{
		AuthHandler:        authHandler,
	})

	return r
}
