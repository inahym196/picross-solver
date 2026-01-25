package accessor

import (
	"fmt"

	"github.com/inahym196/picross-solver/pkg/game"
)

type lineKind uint8

const (
	LineKindRow lineKind = iota
	LineKindColumn
)

func (kind lineKind) String() string {
	switch kind {
	case LineKindRow:
		return "Row"
	case LineKindColumn:
		return "Col"
	default:
		panic("invalid lineKind")
	}
}

// deducer.DeduceLineが参照しているためprivateにできない
type LineRef struct {
	kind  lineKind
	index int
}

func (ref LineRef) String() string {
	return fmt.Sprintf("%s[%d]", ref.kind, ref.index)
}

type LineAccessor struct {
	game game.Game
	ref  LineRef
}

func NewLineAccessor(game game.Game, kind lineKind, index int) LineAccessor {
	return LineAccessor{game, LineRef{kind, index}}
}

func (acc LineAccessor) Cells() []game.Cell {
	board := acc.game.Board()
	index := acc.ref.index

	switch acc.ref.kind {
	case LineKindRow:
		return board.Cells()[index]
	case LineKindColumn:
		cells := make([]game.Cell, board.GetRows())
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
