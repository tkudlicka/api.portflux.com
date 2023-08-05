package entities

import (
	"time"
)

// EntityNameUser contains the name of the entity
const EntityNameUser = "user"

// User struct
type User struct {
	UserID       string    `bson:"_userid,omitempty"`
	Firstname    string    `bson:"firstname"`
	Lastname     string    `bson:"lastname"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}
