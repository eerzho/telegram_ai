package genkit_stub

import (
	"context"
	"fmt"

	"github.com/eerzho/telegram-ai/internal/domain"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) GenerateSummary(
	ctx context.Context,
	language string,
	dialog domain.Dialog,
	onChunk func(chunk string) error,
) error {
	const op = "genkit.Stub.GenerateSummary"

	text := generateStubText(500)

	words := splitWords(text)
	for _, word := range words {
		if err := onChunk(word); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func splitWords(text string) []string {
	var words []string
	var current string

	for _, ch := range text {
		current += string(ch)
		if ch == ' ' || ch == '\n' {
			words = append(words, current)
			current = ""
		}
	}

	if current != "" {
		words = append(words, current)
	}

	return words
}

func generateStubText(size int) string {
	const pattern = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. "

	var result string
	for len(result) < size {
		result += pattern
	}

	return result[:size]
}
