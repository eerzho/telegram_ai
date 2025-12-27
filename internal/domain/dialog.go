package domain

type Dialog struct {
	Owner    User
	Messages []Message
}

type User struct {
	ChatID int64
	Name   string
}

type Message struct {
	Sender User
	Text   string
	Date   int
}
