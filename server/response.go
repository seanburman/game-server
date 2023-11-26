package server

type (
	Message struct {
		Name    string `json:"name,omitempty"`
		Message string `json:"message,omitempty"`
	}
)
