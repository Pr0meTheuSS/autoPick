package models

// User - модель пользователя
type User struct {
	ID        string
	Email     string
	Username  string
	Password  string
	Confirmed bool
	Role      string
	IsBlocked bool
}
