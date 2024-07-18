package models

type UserPayload struct {
	Email string `json:"email" binding:"required"`
}

type UserResponse struct {
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
	Error   string  `json:"error,omitempty"`
}
