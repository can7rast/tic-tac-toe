package application

import (
	"errors"
	"school21/internal/domain"
	"school21/internal/infrastructure/datasource"
)

type gameService struct {
	repo datasource.GameRepository
}

func (g *gameService) GetNextMove(game *domain.Game) (domain.Game, error) {
	isOver, err, _ := g.IsGameOver(game)
	if err != nil {
		return domain.Game{}, err
	}
	if isOver {
		return *game, nil
	}

	newGame := *game
	board := newGame.Board.Board

	row, col := domain.FindBestMove(board)
	if row == -1 || col == -1 {
		return newGame, nil
	}

	board[row][col] = domain.Computer
	newGame.Board.Board = board
	newGame.CurrentPlayer = domain.Player

	return newGame, nil
}
func (g *gameService) ValidateMove(newGame *domain.Game) error {
	oldGame, err := g.repo.Get(newGame.ID)
	if err != nil {
		return err
	}

	if newGame.CurrentPlayer != domain.Player {
		return errors.New("сейчас не ход игрока")
	}

	diff := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if oldGame.Board.Board[i][j] != domain.Empty && oldGame.Board.Board[i][j] != newGame.Board.Board[i][j] {
				return errors.New("Изменилась непустая клетка")
			}
			if oldGame.Board.Board[i][j] == domain.Empty && newGame.Board.Board[i][j] != domain.Empty {
				if newGame.Board.Board[i][j] != domain.Player {
					return errors.New("Игрок ставит только X")
				}
				diff++
			}
		}
	}
	if diff != 1 {
		return errors.New("Игрок делает только 1 ход")
	}
	return nil
}

func (g *gameService) IsGameOver(Game *domain.Game) (bool, error, *int) {
	b := Game.Board.Board

	if score := domain.Evaluate(b); score != 0 {
		winner := 0
		if score == 10 {
			winner = domain.Computer
		} else if score == -10 {
			winner = domain.Player
		}
		return true, nil, &winner
	}

	//проверка на ничью
	draw := true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == domain.Empty {
				draw = false
			}
		}
	}
	if draw {
		return true, nil, nil
	}
	return false, nil, nil
}

func NewGameService(repo datasource.GameRepository) domain.GameService {
	return &gameService{repo: repo}
}
