package models

type LoginCredentials struct {
	Email    string
	Password string
}

type AuthCredentials struct {
	AccessToken  string
	RefreshToken string
}
