package responsegenerate

import "github.com/eerzho/telegram-ai/internal/domain"

func inputToDialog(input Input) domain.Dialog {
	messages := make([]domain.Message, 0, len(input.Messages))
	for _, msg := range input.Messages {
		messages = append(messages, domain.Message{
			Sender: domain.User{
				ChatID:   msg.Sender.ChatID,
				Nickname: msg.Sender.Nickname,
			},
			Text: msg.Text,
			Date: msg.Date,
		})
	}

	return domain.Dialog{
		Owner: domain.User{
			ChatID:   input.Owner.ChatID,
			Nickname: input.Owner.Nickname,
		},
		Messages: messages,
	}
}
