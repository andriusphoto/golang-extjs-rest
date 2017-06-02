package main

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

type jsonReturnArray struct {
	Data    []interface{} `json:"data"`
	Success bool          `json:"success"`
	Length  int           `json:"total"`
}

type jsonReturn struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}
type jsonReturnError struct {
	Errors  interface{} `json:"errors"`
	Success bool        `json:"success"`
}
type filter struct {
	Property string `json:"property"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}
type sorter struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}
