package genkit

import (
	"context"
	"errors"
	"os"

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
	genkit *genkit.Genkit
}

func New(cfg Config) *Client {
	_ = os.Setenv("OPENAI_API_KEY", cfg.Key)
	return &Client{
		genkit: genkit.Init(context.Background(),
			genkit.WithPlugins(&openai.OpenAI{}),
			genkit.WithPromptDir("./prompts"),
		),
	}
}
