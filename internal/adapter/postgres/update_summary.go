package postgres

import (
	"context"
	"fmt"
)

func (db *DB) UpdateSummary(ctx context.Context, ownerID, peerID, text string) error {
	const op = "postgres.DB.UpdateSummary"

	var err error

	stmt, ok := db.stmtCache.Get("update_summary")
	if !ok {
		query := `
			insert into summaries (owner_id, peer_id, text)
			values ($1, $2, $3)
			on conflict (owner_id, peer_id)
			do update set text = excluded.text
		`
		stmt, err = db.db.PreparexContext(ctx, query)
		if err != nil {
			return nil
		}
		db.stmtCache.Put("update_summary", stmt)
	}

	_, err = stmt.ExecContext(ctx, ownerID, peerID, text)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
