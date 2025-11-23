package summary_get

type Input struct {
	OwnerID string `validate:"required,min=1"`
	PeerID  string `validate:"required,min=1"`
}

type Output struct {
	Text string `json:"text"`
}
