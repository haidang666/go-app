package entity

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type User struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	HashedPassword string     `json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	errs := validate.Var(u.Email, "required,email")
	if errs != nil {
		return errs
	}
	errs = validate.Var(u.HashedPassword, "required")
	if errs != nil {
		return errors.New("hashed password is required")
	}
	return nil
}
