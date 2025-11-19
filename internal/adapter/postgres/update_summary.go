package postgres

import (
	"context"
	"fmt"
)

func (db *DB) UpdateSummary(ctx context.Context, chatID, text string) error {
	const op = "postgres.DB.UpdateSummary"

	query := `
		insert into summaries (chat_id, text)
		values ($1, $2)
		on conflict (chat_id)
		do update set text = excluded.text
	`
	_, err := db.db.ExecContext(ctx, query, chatID, text)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
