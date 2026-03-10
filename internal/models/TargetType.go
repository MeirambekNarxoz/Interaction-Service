package models

type TargetType string

const (
	TargetPost    = "post"
	TargetArticle = "article"
	TargetComment = "comment"
)

func IsValidTargetType(t TargetType) bool {
	switch t {
	case TargetPost, TargetArticle, TargetComment:
		return true
	default:
		return false
	}
}
