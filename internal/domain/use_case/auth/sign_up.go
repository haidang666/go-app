package auth

import (
	"context"

	"github.com/haidang666/go-app/internal/domain/entity"
	"github.com/haidang666/go-app/internal/domain/repository"
	"github.com/haidang666/go-app/internal/domain/use_case/auth/dto"
	"golang.org/x/crypto/bcrypt"
)

type SignUpUseCase struct {
	userRepo repository.UserRepository
}

func NewSignUpUseCase(userRepo repository.UserRepository) *SignUpUseCase {
	return &SignUpUseCase{userRepo: userRepo}
}

func (uc *SignUpUseCase) Execute(ctx context.Context, req *dto.SignUpRequestDto) (*entity.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	du := &entity.User{
		Email:          req.Email,
		HashedPassword: string(hashed),
	}

	if err := du.Validate(); err != nil {
		return nil, err
	}
	newUser, err := uc.userRepo.Create(ctx, du)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
