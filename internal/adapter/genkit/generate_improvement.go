package genkit

import (
	"context"
	"fmt"

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
		return fmt.Errorf("%s: %w", op, ErrPromptNotFound)
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
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
