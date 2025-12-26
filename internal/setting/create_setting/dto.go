package createsetting

import "github.com/eerzho/telegram_ai/internal/domain"

type Input struct {
	UserID int64  `json:"user_id" validate:"required"`
	ChatID int64  `json:"chat_id" validate:"required"`
	Style  string `json:"text" validate:"max=500"`
}

type Output struct {
	domain.Setting
}
