package postgres

import (
	"context"
	"fmt"
)

func (db *DB) UpdateSummary(ctx context.Context, chatID, text string) error {
	const op = "postgres.DB.UpdateSummary"

	var err error

	stmt, ok := db.stmtCache.Get("update_summary")
	if !ok {
		query := `
			insert into summaries (chat_id, text)
			values ($1, $2)
			on conflict (chat_id)
			do update set text = excluded.text
		`
		stmt, err = db.db.PreparexContext(ctx, query)
		if err != nil {
			return nil
		}
		db.stmtCache.Put("update_summary", stmt)
	}

	_, err = stmt.ExecContext(ctx, chatID, text)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
