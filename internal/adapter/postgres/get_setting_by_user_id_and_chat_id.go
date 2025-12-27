package postgres

import (
	"context"
	"database/sql"
	"errors"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
)

func (db *DB) GetSettingByUserIDAndChatID(ctx context.Context, userID, chatID int64) (domain.Setting, error) {
	const op = "postgres.DB.GetSettingByUserIDAndChatID"

	query := `
		select user_id, chat_id, style, created_at, updated_at from settings
		where user_id = $1 and chat_id = $2
	`
	var setting domain.Setting
	if err := db.GetContext(ctx, &setting, query, userID, chatID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Setting{}, errorhelp.WithOP(op, domain.ErrSettingNotFound)
		}
		return domain.Setting{}, errorhelp.WithOP(op, err)
	}

	return setting, nil
}
