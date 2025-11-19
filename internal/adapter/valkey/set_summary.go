package valkey

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) SetSummary(ctx context.Context, chatID, text string) error {
	const op = "valkey.Client.SetSummary"

	key := fmt.Sprintf("%s:%s", summaryPrefix, chatID)

	cmd := c.client.Do(ctx, c.client.B().Set().Key(key).Value(text).Ex(time.Hour).Build())
	if err := cmd.Error(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
