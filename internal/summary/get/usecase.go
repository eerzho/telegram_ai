package get

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Cacher interface {
	GetSummary(ctx context.Context, chatID string) (string, error)
}

type Usecase struct {
	validate *validator.Validate
	cacher   Cacher
}

func NewUsecase(
	validate *validator.Validate,
	cacher Cacher,
) *Usecase {
	return &Usecase{
		validate: validate,
		cacher:   cacher,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "summary.get.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	text, err := u.cacher.GetSummary(ctx, input.ChatID)
	if err != nil {
		return Output{}, fmt.Errorf("%s: %w", op, err)
	}

	return Output{Text: text}, nil
}
