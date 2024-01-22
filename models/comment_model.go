package models

type CommentRequest struct {
	MovieID     int64  `json:"movie_id"`
	CommentText string `json:"comment_text"`
}

type CommentResponse struct {
	ID          int64  `json:"id"`
	UserID      string `json:"user_id"`
	MovieID     int64  `json:"movie_id"`
	CommentText string `json:"comment_text"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

type CommentDetails struct {
	Username    string `json:"username"`
	CommentText string `json:"comment_text"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}
