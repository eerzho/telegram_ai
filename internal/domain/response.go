package domain

type ResponseType int

const (
	ResponseTypeMessage ResponseType = iota
	ResponseTypeReaction
	ResponseTypeSkip
)

type ReactionType int

const (
	ReactionTypeZeroVal ReactionType = iota
	ReactionTypeLike
	ReactionTypeOK
	ReactionTypeNice
)

type Response struct {
	Type         ResponseType `json:"type"`
	Message      string       `json:"message"`
	ReactionType ReactionType `json:"reaction_type"`
}
