package valkey

import (
	"context"
	"fmt"
)

func (c *Client) GetSummary(ctx context.Context, chatID string) (string, error) {
	const op = "valkey.Client.GetSummary"

	key := fmt.Sprintf("%s:%s", summaryPrefix, chatID)

	cmd := c.client.Do(ctx, c.client.B().Get().Key(key).Build())
	if err := cmd.Error(); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	text, err := cmd.AsBytes()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(text), nil
}
