package domain

type Message struct {
	Sender User
	Text   string
	Date   int
}

type User struct {
	ChatID string
	Name   string
}

type Dialog struct {
	Owner    User
	Messages []Message
}
