package genkit

import (
	"context"

	errorhelp "github.com/eerzho/telegram-ai/pkg/error_help"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (c *Client) GenerateImprovement(
	ctx context.Context,
	text string,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Client.GenerateImprovement"

	promptName := "improvement"
	prompt := genkit.LookupPrompt(c.genkit, promptName)
	if prompt == nil {
		return errorhelp.WithOP(op, ErrPromptNotFound)
	}

	input := map[string]any{"text": text}
	_, err := prompt.Execute(ctx,
		ai.WithInput(input),
		ai.WithStreaming(func(_ context.Context, aiChunk *ai.ModelResponseChunk) error {
			chunk := aiChunk.Text()
			if chunk != "" {
				return onChunk(chunk)
			}
			return nil
		}),
	)
	if err != nil {
		return errorhelp.WithOP(op, err)
	}
	return nil
}
