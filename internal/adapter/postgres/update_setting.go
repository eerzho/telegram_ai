package postgres

import (
	"context"
	"database/sql"
	"errors"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
)

func (db *DB) UpdateSetting(ctx context.Context, userID, chatID int64, style string) (domain.Setting, error) {
	const op = "postgres.DB.UpdateSetting"

	query := `
		update settings
		set style = $1, updated_at = now()
		where user_id = $2 and chat_id = $3
		returning user_id, chat_id, style, created_at, updated_at
	`
	var setting domain.Setting
	if err := db.GetContext(ctx, &setting, query, style, userID, chatID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Setting{}, errorhelp.WithOP(op, domain.ErrSettingNotFound)
		}
		return domain.Setting{}, errorhelp.WithOP(op, err)
	}

	return setting, nil
}
