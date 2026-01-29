package application

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"school21/internal/domain"
)

type gameService struct {
	repo domain.GameRepository
}

func (g *gameService) ShowAllAvailableGames() ([]domain.Game, error) {
	games, err := g.repo.ShowAllAvailableGames()
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (g *gameService) JoinGame(game *domain.Game, userID uuid.UUID) error {
	if game.State != domain.WaitingPlayers {
		return errors.New("game is not waiting players")
	}
	if game.Player1 == userID {
		return errors.New("you can't join in your game")
	}

	game.Player2 = &userID
	game.State = domain.TurnPlayer1
	return nil
}

func (g *gameService) GetNextMove(game *domain.Game) (domain.Game, error) {
	isOver, err, _ := g.IsGameOver(game)
	if err != nil {
		log.Println(err)
		return domain.Game{}, err
	}
	if isOver {
		log.Println("game is over")
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
	newGame.State = domain.TurnPlayer1

	return newGame, nil
}
func (g *gameService) ValidateMove(newGame *domain.Game, userID uuid.UUID) error {
	oldGame, err := g.repo.Get(newGame.ID)
	if err != nil {
		return err
	}
	newGame.Player1 = oldGame.Player1
	newGame.Player2 = oldGame.Player2
	newGame.VsAI = oldGame.VsAI

	var expectedSign int
	var expectedPlayerUUID uuid.UUID

	switch oldGame.State {
	case domain.TurnPlayer1:
		expectedSign = domain.X
		expectedPlayerUUID = oldGame.Player1
	case domain.TurnPlayer2:
		expectedSign = domain.O
		if oldGame.Player2 != nil {
			expectedPlayerUUID = *oldGame.Player2
		} else {
			newGame.VsAI = true
		}
	default:
		return errors.New("игра не в состоянии, где можно делать ход")
	}

	if expectedPlayerUUID != userID {
		return errors.New("сейчас не ваш ход")
	}

	diff := 0
	var placedSign int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			oldVal := oldGame.Board.Board[i][j]
			newVal := newGame.Board.Board[i][j]
			if oldVal != newVal {
				if oldVal != domain.Empty {
					return errors.New("нельзя перезаписать заполненную клетку")
				}
				diff++
				placedSign = newVal
			}
		}
	}
	if diff != 1 {
		return errors.New("вы должны сделать ровно один ход")
	}
	if expectedSign != placedSign {
		return fmt.Errorf("ожидался символ %d, поставлен символ %d", expectedSign, placedSign)
	}

	if oldGame.State == domain.TurnPlayer1 {
		newGame.State = domain.TurnPlayer2
		newGame.CurrentPlayer = domain.SecondPlayer
	} else if oldGame.State == domain.TurnPlayer2 {
		newGame.State = domain.TurnPlayer1
		newGame.CurrentPlayer = domain.Player
	}

	return nil
}

func (g *gameService) IsGameOver(Game *domain.Game) (bool, error, *int) {
	b := Game.Board.Board

	if score := domain.Evaluate(b); score != 0 {
		winner := 0
		if score == 10 {
			if Game.Player2 != nil {
				winner = domain.SecondPlayer
			} else {
				winner = domain.Computer
			}
			Game.State = domain.WinPlayer2
		} else if score == -10 {
			winner = domain.Player
			Game.State = domain.WinPlayer1
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
		Game.State = domain.Draw
		return true, nil, nil
	}
	return false, nil, nil
}

func NewGameService(repo domain.GameRepository) domain.GameService {
	return &gameService{repo: repo}
}
