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

type GameRepository interface {
	Save(game domain.Game) error
	Get(id uuid.UUID) (domain.Game, error)
}

type gameRepository struct {
	db *DB
}

func withTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, 5*time.Second)
}

func (g *gameRepository) Save(game domain.Game) error {
	boardJson, err := json.Marshal(game.Board)
	if err != nil {
		return fmt.Errorf("could not marshal board: %w", err)
	}

	ctx, cancel := withTimeout(context.Background())
	defer cancel()

	_, err = g.db.Pool.Exec(ctx,
		`INSERT INTO games (id, board, current_player)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (id) DO UPDATE
			 SET board = EXCLUDED.board,
			 current_player = EXCLUDED.current_player`,
		game.ID, boardJson, game.CurrentPlayer)

	return err
}

func (g *gameRepository) Get(id uuid.UUID) (domain.Game, error) {
	var game domain.Game
	var boardJson []byte

	ctx, cancel := withTimeout(context.Background())
	defer cancel()

	err := g.db.Pool.QueryRow(ctx,
		`SELECT id, board, current_player FROM games WHERE id = $1`, id).Scan(&game.ID, &boardJson, &game.CurrentPlayer)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Game{}, fmt.Errorf("Game not found")
		}
		return domain.Game{}, err
	}
	if err = json.Unmarshal(boardJson, &game.Board); err != nil {
		return domain.Game{}, fmt.Errorf("Could not unmarshal board: %w", err)
	}
	return game, nil
}

func NewGameRepository(dataBase *DB) GameRepository {
	return &gameRepository{db: dataBase}
}
