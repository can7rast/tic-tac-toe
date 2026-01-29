package domain

const (
	Empty = iota
	X
	O
)

const (
	Player       = X
	Computer     = O
	SecondPlayer = O
)

type Board struct {
	Board [3][3]int `db:"board"`
}
