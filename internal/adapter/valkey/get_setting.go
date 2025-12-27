package valkey

import (
	"context"
	"encoding/json"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/valkey-io/valkey-go"
)

func (c *Client) GetSetting(ctx context.Context, userID, chatID int64) (domain.Setting, error) {
	const op = "valkey.Client.GetSetting"

	key := c.settingKey(userID, chatID)
	cmd := c.B().Get().Key(key).Build()
	value, err := c.Do(ctx, cmd).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return domain.Setting{}, errorhelp.WithOP(op, domain.ErrSettingNotFound)
		}
		return domain.Setting{}, errorhelp.WithOP(op, err)
	}

	var setting domain.Setting
	if err = json.Unmarshal([]byte(value), &setting); err != nil {
		return domain.Setting{}, errorhelp.WithOP(op, err)
	}

	return setting, nil
}
