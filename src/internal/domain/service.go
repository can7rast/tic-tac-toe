package domain

type GameService interface {
	GetNextMove(game *Game) (Game, error)
	ValidateMove(newGame *Game) error
	IsGameOver(Game *Game) (bool, error, *int)
}
