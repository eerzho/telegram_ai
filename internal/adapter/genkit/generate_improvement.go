package genkit

import (
	"context"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func (c *Client) GenerateImprovement(
	ctx context.Context,
	text string,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Client.GenerateImprovement"

	promptName, data := c.improvementData(text)

	prompt := genkit.LookupPrompt(c.genkit, promptName)
	if prompt == nil {
		return errorhelp.WithOP(op, ErrPromptNotFound)
	}

	_, err := prompt.Execute(ctx,
		ai.WithInput(data),
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

func (c *Client) improvementData(text string) (string, map[string]any) {
	promptName := "improvement"
	data := map[string]any{"text": text}
	return promptName, data
}
