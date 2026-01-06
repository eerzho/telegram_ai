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

func (c *Client) GenerateSummary(
	ctx context.Context,
	language string,
	dialog domain.Dialog,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Client.GenerateSummary"

	promptName := "generate_summary"
	data := c.summaryData(language, dialog)

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

func (c *Client) summaryData(language string, dialog domain.Dialog) map[string]any {
	slices.SortFunc(dialog.Messages, func(a, b domain.Message) int {
		return cmp.Compare(a.Date, b.Date)
	})

	var conversationBuilder strings.Builder

	for _, msg := range dialog.Messages {
		conversationBuilder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			time.Unix(int64(msg.Date), 0).Format(time.DateTime),
			msg.Sender.Name,
			msg.Text,
		))
	}

	data := map[string]any{
		"language":          language,
		"author_name":       dialog.Owner.Name,
		"current_timestamp": time.Now().Format(time.DateTime),
		"conversation":      conversationBuilder.String(),
	}

	return data
}
