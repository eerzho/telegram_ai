package genkit

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type decision struct {
	ResponseType string `json:"response_type"`
	Message      string `json:"message"`
	ReactionType string `json:"reaction_type"`
	Reasoning    string `json:"reasoning"`
}

func (d decision) toResponse() (domain.Response, error) {
	switch d.ResponseType {
	case string(domain.ResponseTypeMessage):
		return domain.Response{
			Type:    domain.ResponseTypeMessage,
			Message: d.Message,
		}, nil
	case string(domain.ResponseTypeReaction):
		switch d.ReactionType {
		case string(domain.ReactionTypeLike):
			return domain.Response{
				Type:         domain.ResponseTypeReaction,
				ReactionType: domain.ReactionTypeLike,
			}, nil
		case string(domain.ReactionTypeOK):
			return domain.Response{
				Type:         domain.ResponseTypeReaction,
				ReactionType: domain.ReactionTypeOK,
			}, nil
		case string(domain.ReactionTypeNice):
			return domain.Response{
				Type:         domain.ResponseTypeReaction,
				ReactionType: domain.ReactionTypeNice,
			}, nil
		default:
			return domain.Response{}, domain.ErrInvalidReactionType
		}
	case string(domain.ResponseTypeSkip):
		return domain.Response{
			Type: domain.ResponseTypeSkip,
		}, nil
	default:
		return domain.Response{}, domain.ErrInvalidReactionType
	}
}

func (c *Client) GenerateResponse(ctx context.Context, userStyle string, dialog domain.Dialog) (domain.Response, error) {
	const op = "genkit.Client.GenerateResponse"

	promptName := "generate_response"
	data := c.responseData(userStyle, dialog)

	prompt := genkit.LookupPrompt(c.genkit, promptName)
	if prompt == nil {
		return domain.Response{}, errorhelp.WithOP(op, ErrPromptNotFound)
	}

	result, err := prompt.Execute(ctx, ai.WithInput(data))
	if err != nil {
		return domain.Response{}, errorhelp.WithOP(op, err)
	}

	var d decision
	if err := json.Unmarshal([]byte(result.Text()), &d); err != nil {
		return domain.Response{}, errorhelp.WithOP(op, err)
	}

	response, err := d.toResponse()
	if err != nil {
		return domain.Response{}, errorhelp.WithOP(op, err)
	}

	return response, nil
}

func (c *Client) responseData(userStyle string, dialog domain.Dialog) map[string]any {
	slices.SortFunc(dialog.Messages, func(a, b domain.Message) int {
		return cmp.Compare(a.Date, b.Date)
	})

	var ownerMessagesBuilder strings.Builder
	var conversationBuilder strings.Builder

	for _, msg := range dialog.Messages {
		conversationBuilder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			time.Unix(int64(msg.Date), 0).Format(time.DateTime),
			msg.Sender.Name,
			msg.Text,
		))

		if msg.Sender.ChatID == dialog.Owner.ChatID {
			ownerMessagesBuilder.WriteString(fmt.Sprintf("%s\n", msg.Text))
		}
	}

	data := map[string]any{
		"user_style":             userStyle,
		"owner_name":             dialog.Owner.Name,
		"owner_message_examples": ownerMessagesBuilder.String(),
		"conversation":           conversationBuilder.String(),
	}

	return data
}
