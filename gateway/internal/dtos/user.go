package dtos

type CreateUserDto struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDto struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Confirmed bool   `json:"confirmed"`
	Role      string `json:"role"`
	IsBlocked bool   `json:"is_blocked"`
}
