package models

type UserModel struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
