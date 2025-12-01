package contract

import (
	"context"

	"github.com/haidang666/go-app/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) (*entity.User, error)
}
