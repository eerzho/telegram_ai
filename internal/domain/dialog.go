package domain

type Dialog struct {
	Language string
	Owner    User
	Messages []Message
}

type User struct {
	ChatID string
	Name   string
}

type Message struct {
	Sender User
	Text   string
	Date   int
}
