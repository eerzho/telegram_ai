package postgres

import (
	"context"

	errorhelp "github.com/eerzho/goiler/pkg/error_help"
	"github.com/eerzho/telegram_ai/internal/domain"
	"github.com/lib/pq"
)

func (db *DB) CreateSetting(ctx context.Context, userID, chatID int64, style string) (domain.Setting, error) {
	const op = "postgres.DB.CreateSetting"

	query := `
		insert into settings(user_id, chat_id, style)
		values($1, $2, $3)
		returning user_id, chat_id, style, created_at, updated_at
	`
	var setting domain.Setting
	if err := db.GetContext(ctx, &setting, query, userID, chatID, style); err != nil {
		if pqErr, ok := errorhelp.AsType[*pq.Error](err); ok && pqErr.Code == "23505" {
			return domain.Setting{}, errorhelp.WithOP(op, domain.ErrSettingAlreadyExists)
		}
		return domain.Setting{}, errorhelp.WithOP(op, err)
	}

	return setting, nil
}
