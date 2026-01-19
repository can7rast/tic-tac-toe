package domain

// evaluate возвращает оценку позиции с точки зрения компьютера:
// +10 — компьютер выиграл
// -10 — игрок выиграл
// 0   — ничья или позиция без немедленного победителя
func Evaluate(board [3][3]int) int {
	b := board

	// Проверка строк
	for i := 0; i < 3; i++ {
		if b[i][0] != Empty && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			if b[i][0] == Computer {
				return 10
			}
			if b[i][0] == Player {
				return -10
			}
		}
	}

	// Проверка столбцов
	for j := 0; j < 3; j++ {
		if b[0][j] != Empty && b[0][j] == b[1][j] && b[1][j] == b[2][j] {
			if b[0][j] == Computer {
				return 10
			}
			if b[0][j] == Player {
				return -10
			}
		}
	}

	// Главная диагональ
	if b[0][0] != Empty && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		if b[0][0] == Computer {
			return 10
		}
		if b[0][0] == Player {
			return -10
		}
	}

	// Побочная диагональ
	if b[0][2] != Empty && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		if b[0][2] == Computer {
			return 10
		}
		if b[0][2] == Player {
			return -10
		}
	}

	return 0 // нет победителя
}

// isBoardFull проверяет, заполнена ли доска полностью (для определения ничьей)
func isBoardFull(board *[3][3]int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if (*board)[i][j] == Empty {
				return false
			}
		}
	}
	return true
}

// Константы для альфа-бета отсечения
const (
	negInf = -999999
	posInf = 999999
)

// minimax — рекурсивная функция с альфа-бета отсечением
// board — указатель, чтобы эффективно менять и откатывать ходы
// depth — текущая глубина (используется для предпочтения быстрых побед)
// alpha, beta — границы отсечения
// maximizing — true, если ходит компьютер (максимизируем оценку)
func minimax(board *[3][3]int, depth int, alpha, beta int, maximizing bool) int {
	score := Evaluate(*board)

	// Терминальные состояния
	if score == 10 {
		return score - depth // чем раньше победа компьютера — тем лучше
	}
	if score == -10 {
		return score + depth // чем позже победа игрока — тем лучше для компьютера
	}
	if isBoardFull(board) {
		return 0 // ничья
	}

	if maximizing { // ход компьютера
		best := negInf
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if (*board)[i][j] == Empty {
					(*board)[i][j] = Computer
					eval := minimax(board, depth+1, alpha, beta, false)
					(*board)[i][j] = Empty

					if eval > best {
						best = eval
					}
					if eval > alpha {
						alpha = eval
					}
					if beta <= alpha {
						return best // бета-отсечение
					}
				}
			}
		}
		return best
	} else { // ход игрока (минимизируем)
		best := posInf
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if (*board)[i][j] == Empty {
					(*board)[i][j] = Player
					eval := minimax(board, depth+1, alpha, beta, true)
					(*board)[i][j] = Empty

					if eval < best {
						best = eval
					}
					if eval < beta {
						beta = eval
					}
					if beta <= alpha {
						return best // альфа-отсечение
					}
				}
			}
		}
		return best
	}
}

// FindBestMove возвращает координаты (row, col) лучшего хода для компьютера.
// Если ходов нет — возвращает (-1, -1)
func FindBestMove(board [3][3]int) (int, int) {
	bestScore := negInf
	bestRow, bestCol := -1, -1

	// Работаем с копией доски, чтобы не портить оригинал
	tempBoard := board

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if tempBoard[i][j] == Empty {
				// Пробуем ход
				tempBoard[i][j] = Computer
				score := minimax(&tempBoard, 0, negInf, posInf, false) // после хода компьютера минимизирует игрок
				tempBoard[i][j] = Empty                                // откат

				if score > bestScore {
					bestScore = score
					bestRow = i
					bestCol = j
				}
			}
		}
	}

	return bestRow, bestCol
}
