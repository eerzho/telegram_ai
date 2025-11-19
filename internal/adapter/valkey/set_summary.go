package valkey

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) SetSummary(ctx context.Context, chatID, text string) error {
	const op = "valkey.Client.SetSummary"

	key := fmt.Sprintf("%s:%s", summaryKeyPrefix, chatID)

	now := time.Now()
	endOfDay := now.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	ttl := endOfDay.Sub(now)

	cmd := c.client.Do(ctx, c.client.B().Set().Key(key).Value(text).Ex(ttl).Build())
	if err := cmd.Error(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
