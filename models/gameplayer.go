package models

type PlayerStatus string

const (
	DEAD  PlayerStatus = "DEAD"
	ALIVE PlayerStatus = "ALIVE"
)

type GamePlayer struct {
	GameId   int          `json:"game_id,omitempty"`
	UserId   string       `json:"user_id,omitempty"`
	KillCode *string      `json:"kill_code,omitempty"`
	TargetId *int         `json:"target_id,omitempty"`
	Status   PlayerStatus `json:"status,omitempty"`
}
