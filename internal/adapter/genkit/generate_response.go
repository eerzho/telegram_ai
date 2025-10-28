package genkit

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (c *Client) GenerateResponse(
	ctx context.Context,
	dialog domain.Dialog,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Client.GenerateResponse"

	promptName, input, err := c.createInputForResponse(dialog)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	g := c.createGenkit(ctx)
	prompt := genkit.LookupPrompt(g, promptName)
	if prompt == nil {
		return fmt.Errorf("%s: %w", op, ErrPromptNotFound)
	}

	_, err = prompt.Execute(ctx,
		ai.WithInput(input),
		ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			text := chunk.Text()
			if text != "" {
				return onChunk(text)
			}
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Client) createInputForResponse(dialog domain.Dialog) (string, map[string]any, error) {
	const (
		maxAuthorMessages      = 50
		maxConversationContext = 20
	)

	promptName := "response_without_style"
	hasAuthorMessages := false

	var authorMessagesBuilder strings.Builder
	var conversationBuilder strings.Builder

	authorMessageCount := 0
	for _, msg := range dialog.Messages {
		if dialog.Owner.ChatID == msg.Sender.ChatID {
			hasAuthorMessages = true
			if authorMessageCount < maxAuthorMessages {
				authorMessagesBuilder.WriteString(fmt.Sprintf("- %s\n", msg.Text))
				authorMessageCount++
			}
		}
	}

	contextMessages := dialog.Messages
	if len(contextMessages) > maxConversationContext {
		contextMessages = contextMessages[len(contextMessages)-maxConversationContext:]
	}
	for _, msg := range contextMessages {
		timestamp := time.Unix(int64(msg.Date), 0).Format(time.DateTime)
		conversationBuilder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			timestamp,
			msg.Sender.Name,
			msg.Text,
		))
	}

	input := map[string]any{
		"author_name":       dialog.Owner.Name,
		"current_timestamp": time.Now().Format(time.DateTime),
		"conversation":      conversationBuilder.String(),
	}
	if hasAuthorMessages {
		promptName = "response_with_style"
		input["author_messages"] = strings.TrimSpace(authorMessagesBuilder.String())
	}

	return promptName, input, nil
}
