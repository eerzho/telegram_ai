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

func (c *Client) GenerateSummary(
	ctx context.Context,
	language string,
	dialog domain.Dialog,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Client.GenerateSummary"

	promptName := "summary"
	input := c.createInputForSummary(language, dialog)
	prompt := genkit.LookupPrompt(c.genkit, promptName)
	if prompt == nil {
		return fmt.Errorf("%s: %w", op, ErrPromptNotFound)
	}

	_, err := prompt.Execute(ctx,
		ai.WithInput(input),
		ai.WithStreaming(func(_ context.Context, chunk *ai.ModelResponseChunk) error {
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

func (c *Client) createInputForSummary(language string, dialog domain.Dialog) map[string]any {
	var conversationBuilder strings.Builder

	for _, msg := range dialog.Messages {
		timestamp := time.Unix(int64(msg.Date), 0).Format(time.DateTime)
		conversationBuilder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			timestamp,
			msg.Sender.Nickname,
			msg.Text,
		))
	}

	input := map[string]any{
		"language":          language,
		"author_name":       dialog.Owner.Nickname,
		"current_timestamp": time.Now().Format(time.DateTime),
		"conversation":      conversationBuilder.String(),
	}

	return input
}
