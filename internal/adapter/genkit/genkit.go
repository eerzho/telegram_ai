package genkit

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
)

var (
	ErrPromptNotFound = errors.New("prompt not found")
)

type Config struct {
	Key string `env:"OPENAI_API_KEY"`
}

type Client struct {
}

func New(cfg Config) *Client {
	_ = os.Setenv("OPENAI_API_KEY", cfg.Key)
	return &Client{}
}

func (c *Client) GenerateResponse(
	ctx context.Context,
	dialog string,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Client.GenerateResponse"

	g := c.createGenkit(ctx)

	prompt := genkit.LookupPrompt(g, "response")
	if prompt == nil {
		return fmt.Errorf("%s: %w", op, ErrPromptNotFound)
	}

	_, err := prompt.Execute(ctx,
		ai.WithInput(map[string]any{"dialog": dialog}),
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

func (c *Client) createGenkit(ctx context.Context) *genkit.Genkit {
	return genkit.Init(ctx,
		genkit.WithPlugins(&openai.OpenAI{}),
		genkit.WithPromptDir("./prompts"),
	)
}
