package generateresponse

import (
	"context"
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
}

type generator interface {
	GenerateResponse(ctx context.Context, userStyle string, dialog domain.Dialog) (domain.Response, error)
}

type Usecase struct {
	logger    *slog.Logger
	validate  *validator.Validate
	storage   storage
	cache     cache
	generator generator
}

func NewUsecase(
	logger *slog.Logger,
	validate *validator.Validate,
	storage storage,
	cache cache,
	generator generator,
) *Usecase {
	return &Usecase{
		logger:    logger,
		validate:  validate,
		storage:   storage,
		cache:     cache,
		generator: generator,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	const op = "generate_response.Usecase.Execute"

	if err := u.validate.Struct(input); err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	userStyle := u.getUserStyle(ctx, input.UserID, input.ChatID)

	response, err := u.generator.GenerateResponse(ctx, userStyle, input.ToDialog())
	if err != nil {
		return Output{}, errorhelp.WithOP(op, err)
	}

	return Output{Response: response}, nil
}

func (u *Usecase) getUserStyle(ctx context.Context, userID, chatID int64) string {
	const op = "generate_response.Usecase.getUserStyle"

	setting, err := u.cache.GetSetting(ctx, userID, chatID)
	if err == nil {
		return setting.Style
	}
	u.logger.InfoContext(ctx, "failed to get setting from cache", slog.Any("error", errorhelp.WithOP(op, err)))

	setting, err = u.storage.GetSettingByUserIDAndChatID(ctx, userID, chatID)
	if err == nil {
		return setting.Style
	}
	u.logger.InfoContext(ctx, "failed to get setting from storage", slog.Any("error", errorhelp.WithOP(op, err)))

	return ""
}
