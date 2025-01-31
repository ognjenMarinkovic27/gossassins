package models

type ApprovalStatus string

const (
	NotApproved ApprovalStatus = "NOT_APPROVED"
	Approved    ApprovalStatus = "APPROVED"
)

type GameApproval struct {
	GameId *int           `json:"game_id,omitempty"`
	UserId string         `json:"user_id,omitempty"`
	Status ApprovalStatus `json:"status,omitempty"`
}
