package generateresponse

import (
	"github.com/eerzho/telegram_ai/internal/domain"
)

type Input struct {
	UserID   int64          `json:"user_id"  validate:"required"                     swaggerignore:"true"`
	ChatID   int64          `json:"chat_id"  validate:"required"                     swaggerignore:"true"`
	Owner    InputUser      `json:"owner"    validate:"required"`
	Messages []InputMessage `json:"messages" validate:"required,min=1,max=1000,dive"`
}

type InputUser struct {
	ChatID int64  `json:"chat_id" validate:"required"`
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
	domain.Response
}
