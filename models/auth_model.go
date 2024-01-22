package models

// Credentials representa las credenciales del usuario
type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorLogin struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
