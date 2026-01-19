package dto

import (
	"github.com/google/uuid"
	"school21/internal/domain"
)

type SignUpRequest struct {
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}
type GameRequest struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Board         [3][3]int `json:"board" db:"board"`
	CurrentPlayer int       `json:"currentPlayer" db:"current_player"`
}

type GameResponse struct {
	ID            uuid.UUID `json:"id"`
	Board         [3][3]int `json:"board"`
	CurrentPlayer int       `json:"currentPlayer"`
	GameOver      bool      `json:"gameOver"`
	Winner        *int      `json:"winner"`
}

func (r *GameRequest) ToDomain() domain.Game {
	return domain.Game{
		ID: r.ID,
		Board: domain.Board{
			Board: r.Board,
		},
		CurrentPlayer: r.CurrentPlayer,
	}
}

func FromDomain(g domain.Game, gameOver bool, winner *int) GameResponse {
	return GameResponse{
		ID:            g.ID,
		Board:         g.Board.Board,
		CurrentPlayer: g.CurrentPlayer,
		GameOver:      gameOver,
		Winner:        winner,
	}
}
