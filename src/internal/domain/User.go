package domain

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
}
