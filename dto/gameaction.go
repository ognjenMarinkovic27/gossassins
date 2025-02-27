package dto

type GameActionStartRequest struct {
	GameId *string `json:"game_id,omitempty"`
}

type GameActionKillRequest struct {
	GameId *string `json:"game_id,omitempty"`

	KillCode *string `json:"kill_code,omitempty"`
}
