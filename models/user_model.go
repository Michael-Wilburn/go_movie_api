package models

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	CreatedDate  string `json:"created_date"`
	UpdatedDate  string `json:"updated_date"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}
