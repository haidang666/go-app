package auth

import (
	"net/http"

	"github.com/haidang666/go-app/internal/domain/use_case/auth"
	"github.com/haidang666/go-app/internal/domain/use_case/auth/dto"
	"github.com/haidang666/go-app/pkg/http/request"
)

type NewAuthHandlerArgs struct {
	SignUpUseCase *auth.SignUpUseCase
}

type AuthHandler struct {
	signUpUseCase *auth.SignUpUseCase
}

func NewAuthHandler(args NewAuthHandlerArgs) *AuthHandler {
	return &AuthHandler{
		signUpUseCase: args.SignUpUseCase,
	}
}

func (h *AuthHandler) SignUp(resWriter http.ResponseWriter, r *http.Request) {
	payload := new(dto.SignUpRequestDto)

	if err := request.FromJSON(r, payload); err != nil {
		request.ToJSON(resWriter, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	user, err := h.signUpUseCase.Execute(r.Context(), payload)
	if err != nil {
		request.ToJSON(resWriter, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	request.ToJSON(resWriter, user, http.StatusCreated)
}
