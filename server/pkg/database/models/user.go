package models

type User struct {
	Username string `json:"username" prop:"username"`
	Password string `json:"password" prop:"password"`
}
