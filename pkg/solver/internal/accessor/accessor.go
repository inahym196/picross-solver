package accessor

import (
	"github.com/inahym196/picross-solver/pkg/game"
)

type LineAccessor struct {
	game game.Game
	ref  game.LineRef
}

func NewLineAccessor(g game.Game, kind game.LineKind, index int) LineAccessor {
	return LineAccessor{g, game.LineRef{Kind: kind, Index: index}}
}

func (acc LineAccessor) Cells() []game.Cell {
	board := acc.game.Board()
	index := acc.ref.Index

	switch acc.ref.Kind {
	case game.LineKindRow:
		return board.Cells()[index]
	case game.LineKindColumn:
		cells := make([]game.Cell, board.Height())
		bcells := board.Cells()
		for i := range len(bcells) {
			cells[i] = bcells[i][index]
		}
		return cells
	default:
		panic("invalid linekind accessor")
	}
}

func (acc LineAccessor) Update(cells []game.Cell) {
	switch acc.ref.Kind {
	case game.LineKindRow:
		row := acc.ref.Index
		for i := range cells {
			acc.game.Mark(row, i, cells[i])
		}
	case game.LineKindColumn:
		col := acc.ref.Index
		for i := range cells {
			acc.game.Mark(i, col, cells[i])
		}
	default:
		panic("invalid linekind accessor")
	}
}

func (acc LineAccessor) Ref() game.LineRef { return acc.ref }
