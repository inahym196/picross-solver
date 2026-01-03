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

type lineAccessor interface {
	get() []Cell
	set(cells []Cell)
	ref() lineRef
}

type rowAccessor struct {
	index int
	board *Board
}

var _ lineAccessor = rowAccessor{}

func (acc rowAccessor) get() []Cell {
	return slices.Clone((*acc.board)[acc.index])
}

func (acc rowAccessor) set(cells []Cell) {
	copy((*acc.board)[acc.index], cells)
}

func (acc rowAccessor) ref() lineRef {
	return lineRef{lineKindRow, acc.index}
}

type colAccessor struct {
	index int
	board *Board
}

var _ lineAccessor = colAccessor{}

func (acc colAccessor) get() []Cell {
	cells := make([]Cell, acc.board.GetRows())
	for i := range *acc.board {
		cells[i] = (*acc.board)[i][acc.index]
	}
	return cells
}

func (acc colAccessor) set(cells []Cell) {
	for i := range cells {
		(*acc.board)[i][acc.index] = cells[i]
	}
}

func (acc colAccessor) ref() lineRef {
	return lineRef{lineKindColumn, acc.index}
}
