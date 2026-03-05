package domain

import "github.com/google/uuid"

type GameRepository interface {
	Save(game Game) error
	Get(id uuid.UUID) (Game, error)
	ShowAllAvailableGames() ([]Game, error)
}
