package domain

import (
	"github.com/google/uuid"
)

type Game struct {
	Board         Board     `db:"board"`
	ID            uuid.UUID `db:"id"`
	CurrentPlayer int       `db:"current_player"`
}

func NewGame() *Game {
	return &Game{
		ID:            uuid.New(),
		CurrentPlayer: Player,
	}
}
