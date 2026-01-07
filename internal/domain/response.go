package domain

type ResponseType string

const (
	ResponseTypeMessage  ResponseType = "message"
	ResponseTypeReaction ResponseType = "reaction"
	ResponseTypeSkip     ResponseType = "skip"
)

type ReactionType string

const (
	ReactionTypeOK      ReactionType = "ok"
	ReactionTypeLike    ReactionType = "like"
	ReactionTypeDislike ReactionType = "dislike"
)

type Response struct {
	Type         ResponseType `json:"type"`
	Message      string       `json:"message"`
	ReactionType ReactionType `json:"reaction_type,omitempty"`
}
