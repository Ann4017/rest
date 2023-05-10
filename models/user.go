package models

type User struct {
	I_id    int    `json:"id"`
	S_name  string `json:"name"`
	S_email string `json:"email"`
}
