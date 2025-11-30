package dto

import "github.com/go-playground/validator/v10"

var validate = validator.New(validator.WithRequiredStructEnabled())

type SignUpRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *SignUpRequestDto) Validate() error {
	errs := validate.Var(req.Email, "required,email")
	if errs != nil {
		return errs
	}
	errs = validate.Var(req.Password, "required,min=5")
	if errs != nil {
		return errs
	}
	return nil
}
