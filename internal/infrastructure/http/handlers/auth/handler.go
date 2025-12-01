package auth

import (
	"net/http"

	"github.com/haidang666/go-app/internal/api/auth"
	"github.com/haidang666/go-app/internal/domain/dto"
	authUseCase "github.com/haidang666/go-app/internal/domain/use_case/auth"
	"github.com/haidang666/go-app/pkg/http/request"
)

type NewAuthHandlerArgs struct {
	SignUpUseCase *authUseCase.SignUpUseCase
}

type AuthHandler struct {
	signUpUseCase *authUseCase.SignUpUseCase
}

func NewAuthHandler(args NewAuthHandlerArgs) *AuthHandler {
	return &AuthHandler{
		signUpUseCase: args.SignUpUseCase,
	}
}

func (h *AuthHandler) SignUp(resWriter http.ResponseWriter, r *http.Request) {
	payload := new(auth.SignUpRequest)

	if err := request.FromJSON(r, payload); err != nil {
		request.ToJSON(resWriter, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		request.ToJSON(resWriter, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Convert API DTO to domain DTO
	input := &dto.SignUpInput{
		Email:    payload.Email,
		Password: payload.Password,
	}

	user, err := h.signUpUseCase.Execute(r.Context(), input)
	if err != nil {
		request.ToJSON(resWriter, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	request.ToJSON(resWriter, user, http.StatusCreated)
}
