package models

type TargetType string

const (
	TargetPost    = "post"
	TargetArticle = "article"
	TargetComment     = "comment"
	TargetApplication = "moderator_application"
)

func IsValidTargetType(t TargetType) bool {
	switch t {
	case TargetPost, TargetArticle, TargetComment, TargetApplication:
		return true
	default:
		return false
	}
}
