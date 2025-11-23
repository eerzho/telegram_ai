package domain

type Summary struct {
	OwnerID string `db:"owner_id"`
	PeerID  string `db:"peer_id"`
	Text    string `db:"text"`
}
