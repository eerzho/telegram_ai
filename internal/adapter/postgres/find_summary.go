package postgres

import (
	"context"
	"fmt"

	"github.com/eerzho/telegram-ai/internal/domain"
)

func (db *DB) FindSummary(ctx context.Context, chatID string) (domain.Summary, error) {
	const op = "postgres.DB.FindSummary"

	var summary domain.Summary
	var err error

	stmt, ok := db.stmtCache.Get("find_summary")
	if !ok {
		query := `
			select chat_id, text from summaries
			where chat_id = $1
		`
		stmt, err = db.db.PreparexContext(ctx, query)
		if err != nil {
			return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
		}
		db.stmtCache.Put("find_summary", stmt)
	}

	err = stmt.GetContext(ctx, &summary, chatID)
	if err != nil {
		return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
	}
	return summary, nil
}
