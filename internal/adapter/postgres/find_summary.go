package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eerzho/telegram-ai/internal/domain"
)

func (db *DB) FindSummary(ctx context.Context, ownerID, peerID string) (domain.Summary, error) {
	const op = "postgres.DB.FindSummary"

	var summary domain.Summary
	var err error

	stmt, ok := db.stmtCache.Get("find_summary")
	if !ok {
		query := `
			select owner_id, peer_id, text from summaries
			where owner_id = $1 and peer_id = $2
		`
		stmt, err = db.db.PreparexContext(ctx, query)
		if err != nil {
			return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
		}
		db.stmtCache.Put("find_summary", stmt)
	}

	err = stmt.GetContext(ctx, &summary, ownerID, peerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Summary{}, fmt.Errorf("%s: %w", op, domain.ErrSummaryNotFound)
		}
		return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
	}
	return summary, nil
}
