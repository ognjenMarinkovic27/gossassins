package dto

type GameActionStartRequest struct {
	GameId *int `json:"game_id,omitempty"`

	/* TODO: Through Auth token */
	CallerUserId *string `json:"caller_id,omitempty"`
}

type GameActionApproveRequest struct {
	GameId *int `json:"game_id,omitempty"`

	/* TODO: Through Auth token */
	CallerUserId *string `json:"caller_id,omitempty"`

	ApproveeUserId *string `json:"approvee_id,omitempty"`
}

type GameActionKillRequest struct {
	GameId *int `json:"game_id,omitempty"`

	/* TODO: Through Auth token */
	KillerUserId *string `json:"killer_id,omitempty"`

	KillCode *string `json:"kill_code,omitempty"`
}
