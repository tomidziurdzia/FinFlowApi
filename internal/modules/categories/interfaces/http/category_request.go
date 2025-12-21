package http

type CategoryRequest struct {
	Name string `json:"name"`
	Type *int   `json:"type"`
}