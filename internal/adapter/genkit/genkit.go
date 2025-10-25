package genkit

import (
	"context"
	"fmt"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
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

func (c *Client) GenerateAnswer(ctx context.Context, dialogContext string) (string, error) {
	const op = "genkit.Client.GenerateAnswer"

	g := c.createGenkit(ctx)

	prompt := genkit.LookupPrompt(g, "answer")
	if prompt == nil {
		return "", fmt.Errorf("%s: prompt 'answer' not found", op)
	}

	result, err := prompt.Execute(ctx, ai.WithInput(dialogContext))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if result == nil {
		return "", fmt.Errorf("%s: empty response from prompt", op)
	}

	answer := result.Text()
	if answer == "" {
		return "", fmt.Errorf("%s: empty text in response", op)
	}

	return answer, nil
}

func (c *Client) createGenkit(ctx context.Context) *genkit.Genkit {
	return genkit.Init(ctx,
		genkit.WithPlugins(&openai.OpenAI{}),
		genkit.WithDefaultModel("openai/gpt-4o"),
		genkit.WithPromptDir("./prompts"),
	)
}
