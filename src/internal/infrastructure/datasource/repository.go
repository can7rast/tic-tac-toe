package datasource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"school21/internal/domain"
	"time"
)

type gameRepository struct {
	db *DB
}

func withTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, 5*time.Second)
}

func (g *gameRepository) Save(game domain.Game) error {
	boardJSON, err := json.Marshal(game.Board)
	if err != nil {
		return fmt.Errorf("could not marshal board: %w", err)
	}

	ctx, cancel := withTimeout(context.Background())
	defer cancel()

	query := `
        INSERT INTO games (id, board, current_player, player1_id, player2_id, state, vsai)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            board          = EXCLUDED.board,
            current_player = EXCLUDED.current_player,
            player1_id     = EXCLUDED.player1_id,
            player2_id     = EXCLUDED.player2_id,
            state          = EXCLUDED.state,
            vsai           = EXCLUDED.vsai
    `

	_, err = g.db.Pool.Exec(ctx, query, game.ID, boardJSON, game.CurrentPlayer, game.Player1, game.Player2, game.State, game.VsAI)
	if err != nil {
		return fmt.Errorf("could not save game %s: %w", game.ID.String(), err)
	}

	return nil
}

func (g *gameRepository) Get(id uuid.UUID) (domain.Game, error) {
	var game domain.Game
	var boardJson []byte

	ctx, cancel := withTimeout(context.Background())
	defer cancel()

	err := g.db.Pool.QueryRow(ctx,
		`SELECT id, board, current_player, player1_id, player2_id, state, vsai FROM games WHERE id = $1`, id).Scan(&game.ID, &boardJson, &game.CurrentPlayer, &game.Player1, &game.Player2, &game.State, &game.VsAI)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_ = fmt.Errorf("game not found")
		}
		return domain.Game{}, err
	}
	if err = json.Unmarshal(boardJson, &game.Board); err != nil {
		return domain.Game{}, fmt.Errorf("could not unmarshal board: %w", err)
	}

	return game, nil
}

func (g *gameRepository) ShowAllAvailableGames() ([]domain.Game, error) {
	ctx, cancel := withTimeout(context.Background())
	defer cancel()

	rows, err := g.db.Pool.Query(ctx,
		`SELECT id, current_player, player1_id, player2_id, state, vsai
			 FROM games
			 WHERE state = $1`, domain.WaitingPlayers)
	if err != nil {
		return nil, fmt.Errorf("could not show all games: %w", err)
	}
	var games []domain.Game
	for rows.Next() {
		var game domain.Game
		err = rows.Scan(&game.ID, &game.CurrentPlayer, &game.Player1, &game.Player2, &game.State, &game.VsAI)
		if err != nil {
			return nil, fmt.Errorf("could not show all games: %w", err)
		}
		games = append(games, game)
	}
	return games, rows.Err()
}

func NewGameRepository(dataBase *DB) domain.GameRepository {
	return &gameRepository{db: dataBase}
}
