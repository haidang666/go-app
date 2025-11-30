package infrastructure

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/haidang666/go-app/internal/domain/entity"
	userRepositoryInterface "github.com/haidang666/go-app/internal/domain/repository"
)

type UserRepository struct {
}

var _ userRepositoryInterface.UserRepository = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{
	}
}

func (r *UserRepository) Create(ctx context.Context, du *entity.User) (*entity.User, error) {
	newUser := &entity.User{
		ID:             uuid.New(),
		Email:          strings.ToLower(du.Email),
		HashedPassword: du.HashedPassword,
	}
	return newUser, nil
}

