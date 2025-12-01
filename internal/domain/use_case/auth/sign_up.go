package auth

import (
	"context"

	"github.com/haidang666/go-app/internal/domain/contract"
	"github.com/haidang666/go-app/internal/domain/dto"
	"github.com/haidang666/go-app/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type SignUpUseCase struct {
	userRepo contract.UserRepository
}

func NewSignUpUseCase(userRepo contract.UserRepository) *SignUpUseCase {
	return &SignUpUseCase{userRepo: userRepo}
}

func (uc *SignUpUseCase) Execute(ctx context.Context, input *dto.SignUpInput) (*entity.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	du := &entity.User{
		Email:          input.Email,
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
