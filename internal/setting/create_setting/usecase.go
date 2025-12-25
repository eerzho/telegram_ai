package createsetting

import (
	"context"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type storage interface {
	CreateSetting(ctx context.Context, userID, chatID int64) (domain.Setting, error)
}

type Usecase struct {
	validate *validator.Validate
	storage  storage
}

func NewUsecase(
	validate *validator.Validate,
	storage storage,
) *Usecase {
	return &Usecase{
		validate: validate,
		storage:  storage,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "create_setting.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	setting, err := u.storage.CreateSetting(ctx, input.UserID, input.ChatID)
	if err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	return Output{Setting: setting}, nil
}
