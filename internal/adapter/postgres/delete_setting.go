package postgres

import (
	"context"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
)

func (db *DB) DeleteSetting(ctx context.Context, userID, chatID int64) error {
	const op = "postgres.DB.DeleteSetting"

	query := `
		delete from settings where user_id = $1 and chat_id = $2
	`
	result, err := db.ExecContext(ctx, query, userID, chatID)
	if err != nil {
		return errorhelp.WithOP(op, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errorhelp.WithOP(op, err)
	}
	if affected == 0 {
		return errorhelp.WithOP(op, domain.ErrSettingNotFound)
	}

	return nil
}
