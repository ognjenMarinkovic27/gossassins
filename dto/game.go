package dto

type CreateGameRequest struct {
	Name      string `json:"name,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
}

type PatchGameRequest struct {
	Name *string `json:"name,omitempty"`
}
