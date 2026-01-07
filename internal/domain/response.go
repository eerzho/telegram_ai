package domain

type ResponseType string

const (
	ResponseTypeMessage  ResponseType = "message"
	ResponseTypeReaction ResponseType = "reaction"
	ResponseTypeSkip     ResponseType = "skip"
)

type ReactionType string

const (
	ReactionTypeLike ReactionType = "like"
	ReactionTypeOK   ReactionType = "ok"
	ReactionTypeNice ReactionType = "nice"
)

type Response struct {
	Type         ResponseType `json:"type"`
	Message      string       `json:"message"`
	ReactionType ReactionType `json:"reaction_type,omitempty"`
}
