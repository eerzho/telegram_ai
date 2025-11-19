package get

type Input struct {
	ChatID string `validate:"required,min=1"`
}

type Output struct {
	Text string `json:"text"`
}
