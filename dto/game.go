package dto

type CreateGameRequest struct {
	Name string `json:"name,omitempty"`
}

type PatchGameRequest struct {
	Name *string `json:"name,omitempty"`
}
