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
	ID            uuid.UUID  `json:"id"`
	Board         [3][3]int  `json:"board"`
	CurrentPlayer int        `json:"currentPlayer"`
	GameOver      bool       `json:"gameOver"`
	Winner        *int       `json:"winner"`
	State         int        `json:"state"`
	Player1       uuid.UUID  `json:"player1"`
	Player2       *uuid.UUID `json:"player2"`
	VsAI          bool       `json:"vsAI"`
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
		State:         int(g.State),
		Player1:       g.Player1,
		Player2:       g.Player2,
		VsAI:          g.VsAI,
	}
}

type GameListResponse struct {
	ID      uuid.UUID `json:"id"`
	State   int       `json:"state"`
	Player1 uuid.UUID `json:"player1"`
	VsAi    bool      `json:"vsAI"`
}

func FromDomainList(games []domain.Game) []GameListResponse {
	var resp []GameListResponse
	for i := 0; i < len(games); i++ {
		resp = append(resp, GameListResponse{
			ID:      games[i].ID,
			State:   int(games[i].State),
			Player1: games[i].Player1,
			VsAi:    games[i].VsAI,
		})
	}
	return resp
}
