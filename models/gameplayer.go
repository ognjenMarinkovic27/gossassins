package models

type PlayerStatus string

const (
	NOT_APPROVED PlayerStatus = "NOT_APPROVED"
	ALIVE        PlayerStatus = "ALIVE"
	DEAD         PlayerStatus = "DEAD"
)

type GamePlayer struct {
	GameId   int          `json:"game_id"`
	UserId   string       `json:"user_id"`
	KillCode *string      `json:"kill_code"`
	TargetId *int         `json:"target_id"`
	Status   PlayerStatus `json:"status"`
}

type GamePlayerPatch struct {
	Status *PlayerStatus `json:"status,omitempty"`
}
