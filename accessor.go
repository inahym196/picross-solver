package picrosssolver

import (
	"fmt"
	"slices"
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

type lineAccessor struct {
	board *Board
	ref   lineRef
}

func (acc lineAccessor) Cells() []Cell {
	switch acc.ref.kind {
	case lineKindRow:
		return slices.Clone((*acc.board)[acc.ref.index])
	case lineKindColumn:
		cells := make([]Cell, acc.board.GetRows())
		for i := range *acc.board {
			cells[i] = (*acc.board)[i][acc.ref.index]
		}
		return cells
	default:
		panic("invalid linekind accessor")
	}
}

func (acc lineAccessor) Update(cells []Cell) {
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
