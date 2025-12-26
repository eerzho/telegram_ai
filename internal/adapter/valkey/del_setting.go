package valkey

import (
	"context"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
)

func (c *Client) DelSetting(ctx context.Context, userID, chatID int64) error {
	const op = "valkey.Client.DelSetting"

	key := c.settingKey(userID, chatID)
	cmd := c.B().Del().Key(key).Build()
	err := c.Do(ctx, cmd).Error()
	if err != nil {
		return errorhelp.WithOP(op, err)
	}

	return nil
}
