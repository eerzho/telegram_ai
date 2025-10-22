package input

type StreamAnswer struct {
	Messages []StreamAnswerMessage `json:"messages"`
}

type StreamAnswerMessage struct {
	Text       string `json:"text"`
	Sender     string `json:"sender"`
	Date       int    `json:"date"`
	IsOutgoing bool   `json:"is_outgoing"`
}
