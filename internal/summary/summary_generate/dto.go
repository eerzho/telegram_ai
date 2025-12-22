package summarygenerate

import (
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/eerzho/telegram_ai/pkg/sse"
)

type Input struct {
	Language string         `json:"language" validate:"required,min=2,max=2"`
	Owner    InputUser      `json:"owner"    validate:"required"`
	Messages []InputMessage `json:"messages" validate:"required,min=1,max=1000,dive"`
}

type InputUser struct {
	ChatID string `json:"chat_id" validate:"required,min=1"`
	Name   string `json:"name"    validate:"required,min=1"`
}

type InputMessage struct {
	Sender InputUser `json:"sender" validate:"required"`
	Text   string    `json:"text"   validate:"required,min=1"`
	Date   int       `json:"date"   validate:"required,gt=0"`
}

func (i Input) ToDialog() domain.Dialog {
	messages := make([]domain.Message, 0, len(i.Messages))
	for _, msg := range i.Messages {
		messages = append(messages, msg.ToMessage())
	}
	return domain.Dialog{
		Owner:    i.Owner.ToUser(),
		Messages: messages,
	}
}

func (i InputUser) ToUser() domain.User {
	return domain.User{
		ChatID: i.ChatID,
		Name:   i.Name,
	}
}

func (i InputMessage) ToMessage() domain.Message {
	return domain.Message{
		Sender: i.Sender.ToUser(),
		Text:   i.Text,
		Date:   i.Date,
	}
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
