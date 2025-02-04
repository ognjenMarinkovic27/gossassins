package models

type GameState string

const (
	OPEN    GameState = "OPEN"
	RUNNING GameState = "RUNNING"
	DONE    GameState = "DONE"
)

type Game struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	State     GameState `json:"state"`
}
