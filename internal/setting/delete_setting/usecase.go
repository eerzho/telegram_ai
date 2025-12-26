package deletesetting

import (
	"context"
	"log/slog"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/go-playground/validator/v10"
)

type storage interface {
	DeleteSetting(ctx context.Context, userID, chatID int64) error
}

type cache interface {
	DelSetting(ctx context.Context, userID, chatID int64) error
}

type Usecase struct {
	logger   *slog.Logger
	validate *validator.Validate
	storage  storage
	cache    cache
}

func NewUsecase(
	logger *slog.Logger,
	validate *validator.Validate,
	storage storage,
	cache cache,
) *Usecase {
	return &Usecase{
		logger:   logger,
		validate: validate,
		storage:  storage,
		cache:    cache,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "delete_setting.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	if err := u.storage.DeleteSetting(ctx, input.UserID, input.ChatID); err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	if err := u.cache.DelSetting(ctx, input.UserID, input.ChatID); err != nil {
		u.logger.WarnContext(ctx, "failed to del setting", slog.Any("error", errorhelp.WithOP(op, err)))
	}

	return Output{}, nil
}
