package tests_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"school21/internal/application"
	"school21/internal/domain"
	"school21/internal/domain/mocks"
	"testing"
)

func TestGameService_ShowAllAvailableGames(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)

	mockRepo.On("ShowAllAvailableGames").Return([]domain.Game{
		domain.Game{
			Board:         domain.Board{},
			ID:            uuid.Nil,
			CurrentPlayer: 1,
			State:         domain.TurnPlayer1,
			Player1:       uuid.Nil,
			Player2:       &uuid.Nil,
			VsAI:          false,
		},
	}, nil).Once()

	games, err := service.ShowAllAvailableGames()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(games))
	mockRepo.AssertCalled(t, "ShowAllAvailableGames")
}

func TestGameService_IsGameOver_WinPlayer1(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)

	game := domain.Game{
		Board:         domain.Board{Board: [3][3]int{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}}},
		CurrentPlayer: 1,
		State:         domain.TurnPlayer1,
	}

	isOver, err, win := service.IsGameOver(&game)
	assert.Nil(t, err)
	assert.True(t, isOver)
	assert.Equal(t, domain.Player, *win)
}

func TestGameService_IsGameOver_WinPlayer2(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)

	game := domain.Game{
		Board:         domain.Board{Board: [3][3]int{{2, 0, 0}, {0, 2, 0}, {0, 0, 2}}},
		CurrentPlayer: 1,
		State:         domain.TurnPlayer1,
	}

	isOver, err, win := service.IsGameOver(&game)
	assert.Nil(t, err)
	assert.True(t, isOver)
	assert.Equal(t, domain.SecondPlayer, *win)
}

func TestGameService_GetNextMoveStopLost(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)
	game := domain.Game{
		Board:         domain.Board{Board: [3][3]int{{1, 0, 1}, {0, 0, 0}, {0, 0, 0}}},
		CurrentPlayer: 2,
		State:         domain.TurnPlayer2,
	}

	newGame, err := service.GetNextMove(&game)
	assert.Nil(t, err)
	assert.NotNil(t, newGame)
	assert.Equal(t, 2, newGame.Board.Board[0][1])
}

func TestGameService_GetNext_MoveWin(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)
	game := domain.Game{
		Board:         domain.Board{Board: [3][3]int{{2, 0, 2}, {0, 0, 0}, {0, 0, 0}}},
		CurrentPlayer: 2,
		State:         domain.TurnPlayer2,
	}
	newGame, err := service.GetNextMove(&game)
	assert.Nil(t, err)
	assert.NotNil(t, newGame)
	assert.Equal(t, 2, newGame.Board.Board[0][1])
}

func TestGameService_ValidateMove(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)
	game := domain.Game{
		Board:         domain.Board{Board: [3][3]int{{2, 0, 2}, {0, 0, 0}, {0, 0, 0}}},
		CurrentPlayer: 2,
		State:         domain.TurnPlayer2,
		Player2:       nil,
	}

	mockRepo.EXPECT().Get(mock.Anything).
		Return(game, nil).Once()

	newGame, err := service.GetNextMove(&game)
	assert.Nil(t, err)
	assert.NotNil(t, newGame)
	assert.Equal(t, 2, newGame.Board.Board[0][1])

	err = service.ValidateMove(&newGame, uuid.Nil)
	assert.Nil(t, err)
	assert.Equal(t, true, newGame.VsAI)
}

func TestGameService_ValidateMove_InvalidGame(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)
	game := domain.Game{
		State: domain.WaitingPlayers,
	}

	mockRepo.EXPECT().Get(mock.Anything).
		Return(game, nil).Once()
	newGame, err := service.GetNextMove(&game)
	assert.Nil(t, err)

	err = service.ValidateMove(&newGame, uuid.Nil)
	assert.EqualError(t, err, "игра не в состоянии, где можно делать ход")
}

func TestGameService_ValidateMove_TwoDiff(t *testing.T) {
	mockRepo := mocks.NewGameRepository(t)
	service := application.NewGameService(mockRepo)
	id := uuid.New()
	game := domain.Game{
		Board:         domain.Board{Board: [3][3]int{{1, 1, 1}, {0, 1, 0}, {0, 0, 0}}},
		CurrentPlayer: 1,
		State:         domain.TurnPlayer1,
		Player2:       nil,
		Player1:       id,
	}

	mockRepo.EXPECT().
		Get(mock.Anything).
		Return(domain.Game{
			Board:         domain.Board{Board: [3][3]int{{1, 0, 1}, {0, 0, 0}, {0, 0, 0}}},
			CurrentPlayer: 1,
			State:         domain.TurnPlayer1,
			Player2:       nil,
			Player1:       id,
		}, nil).Once()

	err := service.ValidateMove(&game, id)
	assert.NotNil(t, game)
	assert.EqualError(t, err, "вы должны сделать ровно один ход")

}
