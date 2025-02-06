package dto

type CreateGamePlayerRequest struct {
	UserId *string `json:"user_id,omitempty"`
}

type PatchGamePlayerRequest struct {
	Status *string `json:"status,omitempty"`
}
