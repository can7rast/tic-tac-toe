package domain

import (
	"github.com/google/uuid"
)

type GameService interface {
	GetNextMove(game *Game) (Game, error)
	ValidateMove(newGame *Game, userID uuid.UUID) error
	IsGameOver(Game *Game) (bool, error, *int)
	JoinGame(game *Game, userID uuid.UUID) error
	ShowAllAvailableGames() ([]Game, error)
}
