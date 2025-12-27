package getsetting

import (
	"context"
	"errors"
	"log/slog"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/go-playground/validator/v10"
)

type storage interface {
	GetSettingByUserIDAndChatID(ctx context.Context, userID, chatID int64) (domain.Setting, error)
}

type cache interface {
	GetSetting(ctx context.Context, userID, chatID int64) (domain.Setting, error)
	SetSetting(ctx context.Context, setting domain.Setting) error
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
	const op = "get_setting.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	setting, err := u.cache.GetSetting(ctx, input.UserID, input.ChatID)
	if err == nil {
		return Output{Setting: setting}, nil
	}

	var level slog.Level
	if errors.Is(err, domain.ErrSettingNotFound) {
		level = slog.LevelInfo
	} else {
		level = slog.LevelWarn
	}
	u.logger.Log(ctx, level, "failed to get setting", slog.Any("error", errorhelp.WithOP(op, err)))

	setting, err = u.storage.GetSettingByUserIDAndChatID(ctx, input.UserID, input.ChatID)
	if err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	if err := u.cache.SetSetting(ctx, setting); err != nil {
		u.logger.WarnContext(ctx, "failed to set setting", slog.Any("error", errorhelp.WithOP(op, err)))
	}

	return Output{Setting: setting}, nil
}
