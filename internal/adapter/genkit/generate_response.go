package genkit

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (c *Client) GenerateResponse(
	ctx context.Context,
	dialog domain.Dialog,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Client.GenerateResponse"

	promptName, data := c.responseData(dialog)

	prompt := genkit.LookupPrompt(c.genkit, promptName)
	if prompt == nil {
		return errorhelp.WithOP(op, ErrPromptNotFound)
	}

	_, err := prompt.Execute(ctx,
		ai.WithInput(data),
		ai.WithStreaming(func(_ context.Context, chunk *ai.ModelResponseChunk) error {
			text := chunk.Text()
			if text != "" {
				return onChunk(text)
			}
			return nil
		}),
	)
	if err != nil {
		return errorhelp.WithOP(op, err)
	}

	return nil
}

func (c *Client) responseData(dialog domain.Dialog) (string, map[string]any) {
	slices.SortFunc(dialog.Messages, func(a, b domain.Message) int {
		return cmp.Compare(a.Date, b.Date)
	})

	hasAuthorMessages := false

	var authorMessagesBuilder strings.Builder
	var conversationBuilder strings.Builder

	for _, msg := range dialog.Messages {
		if dialog.Owner.ChatID == msg.Sender.ChatID {
			hasAuthorMessages = true
			authorMessagesBuilder.WriteString(fmt.Sprintf("- %s\n", msg.Text))
		}
		conversationBuilder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			time.Unix(int64(msg.Date), 0).Format(time.DateTime),
			msg.Sender.Name,
			msg.Text,
		))
	}

	promptName := "response_without_style"
	data := map[string]any{
		"author_name":  dialog.Owner.Name,
		"conversation": conversationBuilder.String(),
	}
	if hasAuthorMessages {
		promptName = "response_with_style"
		data["author_messages"] = authorMessagesBuilder.String()
	}

	return promptName, data
}
