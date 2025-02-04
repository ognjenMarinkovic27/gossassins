package models

type JoinRequestStatus string

const (
	NotApproved JoinRequestStatus = "NOT_APPROVED"
	Approved    JoinRequestStatus = "APPROVED"
)

func IsValidJoinRequestStatus(value string) bool {
	switch JoinRequestStatus(value) {
	case Approved, NotApproved:
		return true
	default:
		return false
	}
}

type JoinRequest struct {
	GameId int               `json:"game_id"`
	UserId string            `json:"user_id"`
	Status JoinRequestStatus `json:"status"`
}
