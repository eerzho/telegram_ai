package improvement_generate

type Input struct {
	Text string `json:"text" validate:"required,min=1"`
}

type Output struct {
	TextChan <-chan string
	ErrChan  <-chan error
}
