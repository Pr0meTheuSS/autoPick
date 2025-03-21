package dtos

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthCredentialsDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
