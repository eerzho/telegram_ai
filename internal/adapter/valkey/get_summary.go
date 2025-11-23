package valkey

import (
	"context"
	"fmt"

	"github.com/eerzho/telegram-ai/internal/domain"
	"github.com/valkey-io/valkey-go"
)

func (c *Client) GetSummary(ctx context.Context, ownerID, peerID string) (string, error) {
	const op = "valkey.Client.GetSummary"

	key := c.summaryKey(ownerID, peerID)
	cmd := c.client.Do(ctx, c.client.B().Get().Key(key).Build())
	text, err := cmd.ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", fmt.Errorf("%s: %w", op, domain.ErrSummaryNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return text, nil
}
