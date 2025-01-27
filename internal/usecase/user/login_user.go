package user

import (
	"context"

	"github.com/danielllmuniz/go-http-server/internal/validator"
)

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.NotBlank(req.Password), "password", "this field cannot be empty")

	return eval
}
