package accessor

import (
	"github.com/inahym196/picross-solver/pkg/game"
)

// deducer.DeduceLineが参照しているためprivateにできない

type LineAccessor struct {
	game *game.Game
	ref  LineRef
}

func NewLineAccessor(game *game.Game, kind lineKind, index int) LineAccessor {
	return LineAccessor{game, LineRef{kind, index}}
}

func (acc LineAccessor) Cells() []game.Cell {
	board := acc.game.Board()
	index := acc.ref.index

	switch acc.ref.kind {
	case LineKindRow:
		return board[index]
	case LineKindColumn:
		cells := make([]game.Cell, board.GetRows())
		for i := range board {
			cells[i] = board[i][index]
		}
		return cells
	default:
		panic("invalid linekind accessor")
	}
}

func (acc LineAccessor) Update(cells []game.Cell) {
	switch acc.ref.kind {
	case LineKindRow:
		row := acc.ref.index
		for i := range cells {
			acc.game.SetCell(row, i, cells[i])
		}
	case LineKindColumn:
		col := acc.ref.index
		for i := range cells {
			acc.game.SetCell(i, col, cells[i])
		}
	default:
		panic("invalid linekind accessor")
	}
}

func (acc LineAccessor) Ref() LineRef { return acc.ref }
