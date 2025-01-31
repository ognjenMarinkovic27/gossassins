package models

type PlayerStatus string

const (
	DEAD  PlayerStatus = "DEAD"
	ALIVE PlayerStatus = "ALIVE"
)

type GamePlayer struct {
	GameId   int
	UserId   int
	KillCode string
	TargetId int
	Status   PlayerStatus
}
