package domain

type Summary struct {
	ChatID string `db:"chat_id"`
	Text   string `db:"text"`
}
