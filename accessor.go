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
	game game.Game
	ref  lineRef
}

func (acc lineAccessor) Cells() []game.Cell {
	board := acc.game.Board()
	index := acc.ref.index

	switch acc.ref.kind {
	case lineKindRow:
		return board[index]
	case lineKindColumn:
		cells := make([]game.Cell, board.GetRows())
		for i := range board {
			cells[i] = board[i][index]
		}
		return cells
	default:
		panic("invalid linekind accessor")
	}
}

func (acc lineAccessor) Update(cells []game.Cell) {
	switch acc.ref.kind {
	case lineKindRow:
		row := acc.ref.index
		for i := range cells {
			acc.game.SetCell(row, i, cells[i])
		}
	case lineKindColumn:
		col := acc.ref.index
		for i := range cells {
			acc.game.SetCell(i, col, cells[i])
		}
	default:
		panic("invalid linekind accessor")
	}
}

func (acc lineAccessor) Ref() lineRef { return acc.ref }
