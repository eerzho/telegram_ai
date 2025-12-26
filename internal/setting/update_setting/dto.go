package updatesetting

import "github.com/eerzho/telegram_ai/internal/domain"

type Input struct {
	UserID int64  `json:"user_id" validate:"required" swaggerignore:"true"`
	ChatID int64  `json:"chat_id" validate:"required" swaggerignore:"true"`
	Style  string `json:"style" validate:"max=500"`
}

type Output struct {
	domain.Setting
}
