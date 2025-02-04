package dto

type CreateJoinRequestRequest struct {
	UserId *string `json:"user_id,omitempty"`
}

type PatchJoinRequestRequest struct {
	Status *string `json:"status,omitempty"`
}
