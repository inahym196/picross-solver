package picrosssolver

import (
	"fmt"
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

type lineKind uint8

const (
	lineKindRow lineKind = iota
	lineKindColumn
)

func (kind lineKind) String() string {
	switch kind {
	case lineKindRow:
		return "Row"
	case lineKindColumn:
		return "Col"
	default:
		panic("invalid lineKind")
	}
}

type lineRef struct {
	kind  lineKind
	index int
}

func (ref lineRef) String() string {
	return fmt.Sprintf("%s[%d]", ref.kind, ref.index)
}

type lineView struct {
	Cells []game.Cell
	Hints []int
}

func (line lineView) IsFilled() bool {
	return slices.Index(line.Cells, game.CellUndetermined) == -1
}

type lineAccessor struct {
	board *game.Board
	ref   lineRef
}

func (acc lineAccessor) Cells() []game.Cell {
	switch acc.ref.kind {
	case lineKindRow:
		return slices.Clone((*acc.board)[acc.ref.index])
	case lineKindColumn:
		cells := make([]game.Cell, acc.board.GetRows())
		for i := range *acc.board {
			cells[i] = (*acc.board)[i][acc.ref.index]
		}
		return cells
	default:
		panic("invalid linekind accessor")
	}
}

func (acc lineAccessor) Update(cells []game.Cell) {
	switch acc.ref.kind {
	case lineKindRow:
		copy((*acc.board)[acc.ref.index], cells)
	case lineKindColumn:
		for i := range cells {
			(*acc.board)[i][acc.ref.index] = cells[i]
		}
	default:
		panic("invalid linekind accessor")
	}
}

func (acc lineAccessor) Ref() lineRef { return acc.ref }
