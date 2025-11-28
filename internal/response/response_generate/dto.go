package responsegenerate

import (
	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/eerzho/telegram-ai/pkg/sse"
)

type Input struct {
	Owner    InputUser      `json:"owner"    validate:"required"`
	Messages []InputMessage `json:"messages" validate:"required,min=1,max=1000,dive"`
}

type InputUser struct {
	ChatID   string `json:"chat_id" validate:"required,min=1"`
	Nickname string `json:"name"    validate:"required,min=1"`
}

type InputMessage struct {
	Sender InputUser `json:"sender" validate:"required"`
	Text   string    `json:"text"   validate:"required,min=1"`
	Date   int       `json:"date"   validate:"required,gt=0"`
}

type Output struct {
	TextChan <-chan string
	ErrChan  <-chan error
}

func (i Input) ToDialog() domain.Dialog {
	messages := make([]domain.Message, 0, len(i.Messages))
	for _, msg := range i.Messages {
		messages = append(messages, msg.ToMessage())
	}
	return domain.Dialog{
		Owner: domain.User{
			ChatID:   i.Owner.ChatID,
			Nickname: i.Owner.Nickname,
		},
		Messages: messages,
	}
}

func (i InputMessage) ToMessage() domain.Message {
	return domain.Message{
		Sender: domain.User{
			ChatID:   i.Sender.ChatID,
			Nickname: i.Sender.Nickname,
		},
		Text: i.Text,
		Date: i.Date,
	}
}

func (o *Output) Next() (sse.Event, bool) {
	select {
	case text, ok := <-o.TextChan:
		if !ok {
			return sse.Event{Name: "stop"}, true
		}
		return sse.Event{Name: "append", Data: text}, false
	case err := <-o.ErrChan:
		if err != nil {
			return sse.Event{Name: "stop_with_error", Data: "Please try again later."}, true
		}
		return sse.Event{}, true
	}
}
