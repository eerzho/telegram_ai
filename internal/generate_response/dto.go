package generate_response

type Input struct {
	Owner    InputSender    `json:"owner" validate:"required"`
	Messages []InputMessage `json:"messages" validate:"required,min=1,max=50"`
}

type InputSender struct {
	ChatID int    `json:"chat_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type InputMessage struct {
	Sender InputSender `json:"sender" validate:"required"`
	Text   string      `json:"text" validate:"required"`
	Date   int         `json:"date" validate:"required,gt=0"`
}

type Output struct {
	TextChan <-chan string
	ErrChan  <-chan error
}
