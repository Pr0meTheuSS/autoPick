package models

// User - модель пользователя для MongoDB
type User struct {
	ID        string `bson:"_id,omitempty"` // Автоматически создаваемый ObjectID
	Email     string `bson:"email"`
	Username  string `bson:"username"`
	Password  string `bson:"password"`
	Confirmed bool   `bson:"confirmed"`
	Role      string `bson:"role"`
	IsBlocked bool   `bson:"is_blocked"`
}
