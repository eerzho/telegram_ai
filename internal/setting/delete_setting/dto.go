package deletesetting

type Input struct {
	UserID int64 `json:"user_id" validate:"required"`
	ChatID int64 `json:"chat_id" validate:"required"`
}

type Output struct{}
