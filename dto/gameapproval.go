package dto

type CreateGameApprovalRequest struct {
	UserId *string `json:"user_id,omitempty"`
}

type PatchGameApprovalRequest struct {
	Status *string `json:"status,omitempty"`
}
