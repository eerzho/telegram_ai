package improvementgenerate

import "github.com/eerzho/telegram_ai/pkg/sse"

type Input struct {
	Text string `json:"text" validate:"required,min=1"`
}

type Output struct {
	TextChan <-chan string
	ErrChan  <-chan error
}

func (o *Output) Next() (sse.Event, bool) {
	select {
	case text, ok := <-o.TextChan:
		if !ok {
			return sse.Event{Name: "stop"}, true
		}
		return sse.Event{Name: "append", Data: text}, false
	case err, ok := <-o.ErrChan:
		if !ok {
			return sse.Event{Name: "stop"}, true
		}
		if err != nil {
			return sse.Event{Name: "stop_with_error", Data: "{ \"text\": \"Please try again later.\" }"}, true
		}
		return sse.Event{}, true
	}
}
