package models

type PlayerStatus string

const (
	DEAD  PlayerStatus = "DEAD"
	ALIVE PlayerStatus = "ALIVE"
)

type GamePlayer struct {
	GameId   int          `json:"game_id"`
	UserId   string       `json:"user_id"`
	KillCode *string      `json:"kill_code"`
	TargetId *int         `json:"target_id"`
	Status   PlayerStatus `json:"status"`
}
