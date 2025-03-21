package models

// User - модель пользователя для MongoDB
type RefreshToken struct {
	ID           string `bson:"_id,omitempty"` // Автоматически создаваемый ObjectID
	UserID       string `bson:"user_id"`
	RefreshToken string `bson:"refresh_token"`
}
