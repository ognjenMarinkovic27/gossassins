package models

type ApprovalStatus string

const (
	NotApproved ApprovalStatus = "NOT_APPROVED"
	Approved    ApprovalStatus = "APPROVED"
)

func IsValidApprovalStatus(value string) bool {
	switch ApprovalStatus(value) {
	case Approved, NotApproved:
		return true
	default:
		return false
	}
}

type GameApproval struct {
	GameId int            `json:"game_id"`
	UserId string         `json:"user_id"`
	Status ApprovalStatus `json:"status"`
}
