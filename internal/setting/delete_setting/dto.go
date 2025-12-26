package deletesetting

type Input struct {
	UserID int64 `json:"user_id" validate:"required" swaggerignore:"true"`
	ChatID int64 `json:"chat_id" validate:"required" swaggerignore:"true"`
}

type Output struct{}
