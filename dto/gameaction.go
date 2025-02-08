package dto

type GameActionStartRequest struct {
	GameId *int `json:"game_id,omitempty"`
}

type GameActionKillRequest struct {
	GameId *int `json:"game_id,omitempty"`

	KillCode *string `json:"kill_code,omitempty"`
}
