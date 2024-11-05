package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string `json:"username"`
	Password string `json:"password"`
}
