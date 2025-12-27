package valkey

import (
	"context"
	"encoding/json"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
)

func (c *Client) SetSetting(ctx context.Context, setting domain.Setting) error {
	const op = "valkey.Client.SetSetting"

	key := c.settingKey(setting.UserID, setting.ChatID)
	value, err := json.Marshal(setting)
	if err != nil {
		return errorhelp.WithOP(op, err)
	}

	cmd := c.B().Set().Key(key).Value(string(value)).Ex(c.ttl).Build()
	if err = c.Do(ctx, cmd).Error(); err != nil {
		return errorhelp.WithOP(op, err)
	}

	return nil
}
