package postgres

import (
	"context"
	"fmt"

	"github.com/eerzho/telegram-ai/internal/domain"
)

func (db *DB) FindSummary(ctx context.Context, chatID string) (domain.Summary, error) {
	const op = "postgres.DB.FindSummary"

	var summary domain.Summary
	query := `
		select * from summaries
		where id = $1
	`
	err := db.db.GetContext(ctx, &summary, query, chatID)
	if err != nil {
		return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
	}
	return summary, nil
}
