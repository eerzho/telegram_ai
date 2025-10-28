package generate_summary

type Input struct {
	Owner    InputSender    `json:"owner" validate:"required"`
	Messages []InputMessage `json:"messages" validate:"required,min=1,max=50,dive"`
}

type InputSender struct {
	ChatID string `json:"chat_id" validate:"required,min=1"`
	Name   string `json:"name" validate:"required,min=1"`
}

type InputMessage struct {
	Sender InputSender `json:"sender" validate:"required"`
	Text   string      `json:"text" validate:"required,min=1"`
	Date   int         `json:"date" validate:"required,gt=0"`
}

type Output struct {
	TextChan <-chan string
	ErrChan  <-chan error
}
