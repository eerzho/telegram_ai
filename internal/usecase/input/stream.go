package input

type StreamAnswer struct {
	Owner    StreamAnswerSender    `json:"owner" validate:"required"`
	Messages []StreamAnswerMessage `json:"messages" validate:"required,min=1"`
}

type StreamAnswerSender struct {
	ChatID int    `json:"chat_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type StreamAnswerMessage struct {
	Sender StreamAnswerSender `json:"sender" validate:"required"`
	Text   string             `json:"text" validate:"required"`
	Date   int                `json:"date" validate:"required,gt=0"`
}
