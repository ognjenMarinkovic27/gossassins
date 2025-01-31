package models

type GameState string

const (
	OPEN    GameState = "OPEN"
	CLOSED  GameState = "CLOSED"
	RUNNING GameState = "RUNNING"
	DONE    GameState = "DONE"
)

type Game struct {
	Id        *int      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"`
	State     GameState `json:"state,omitempty"`
}
