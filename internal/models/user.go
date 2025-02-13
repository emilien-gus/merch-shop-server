package models

type User struct {
	ID       int    `json:"-" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"` // Не отправлять пароль в JSON
	Balance  int    `json:"balance" db:"balance"`
}
