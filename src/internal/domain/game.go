package domain

import (
	"github.com/google/uuid"
)

type GameState int

const (
	WaitingPlayers GameState = iota
	TurnPlayer1
	TurnPlayer2
	Draw
	WinPlayer1
	WinPlayer2
)

type Game struct {
	Board         Board      `db:"board"`
	ID            uuid.UUID  `db:"id"`
	CurrentPlayer int        `db:"current_player"`
	State         GameState  `db:"state"`
	Player1       uuid.UUID  `db:"player1"`
	Player2       *uuid.UUID `db:"player2"`
	VsAI          bool       `db:"vs_ai"`
}

func NewGame(creator uuid.UUID, vsAI bool) *Game {
	g := &Game{
		ID:            uuid.New(),
		CurrentPlayer: Player,
		State:         WaitingPlayers,
		Player1:       creator,
		VsAI:          vsAI,
		Player2:       nil,
	}
	if vsAI {
		g.State = TurnPlayer1
	} else {
		g.State = WaitingPlayers
	}
	return g
}
