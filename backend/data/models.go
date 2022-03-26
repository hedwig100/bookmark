package data

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DbUser struct {
	UserId   string
	Username string
	Password string
}
