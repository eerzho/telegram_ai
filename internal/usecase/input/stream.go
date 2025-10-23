package input

type StreamAnswer struct {
	Messages []StreamAnswerMessage `json:"messages" validate:"required,min=1"`
}

type StreamAnswerMessage struct {
	Text       string `json:"text" validate:"required"`
	Sender     string `json:"sender" validate:"required"`
	Date       int    `json:"date" validate:"required,gt=0"`
	IsOutgoing bool   `json:"is_outgoing"`
}
