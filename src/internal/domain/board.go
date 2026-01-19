package domain

const (
	// Empty Пустое поле
	Empty = iota

	// Player Игрок X
	Player

	// Computer  Компьютер O
	Computer
)

type Board struct {
	Board [3][3]int `db:"board"`
}
