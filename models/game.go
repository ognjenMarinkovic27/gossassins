package models

type GameState string

const (
	OPEN    GameState = "OPEN"
	RUNNING GameState = "RUNNING"
	DONE    GameState = "DONE"
)

type Game struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	State     GameState `json:"state"`
	JoinCode  string    `json:"join_code"`
}

type GameWithJoinStatus struct {
	Game
	Joined bool `json:"joined"`
}

type GameCreation struct {
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	State     GameState `json:"state"`
	JoinCode  string    `json:"join_code"`
}

type GamePatch struct {
	Name *string `json:"name,omitempty"`
}
