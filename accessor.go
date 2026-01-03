package picrosssolver

import "slices"

type lineKind uint8

const (
	lineKindRow lineKind = iota
	lineKindColumn
)

type lineRef struct {
	kind  lineKind
	index int
}

type lineAccessor interface {
	Get() []Cell
	Set(cells []Cell)
	Ref() lineRef
}

type rowAccessor struct {
	index int
	board *Board
}

func (acc rowAccessor) Get() []Cell {
	return slices.Clone((*acc.board)[acc.index])
}

func (acc rowAccessor) Set(cells []Cell) {
	copy((*acc.board)[acc.index], cells)
}

func (acc rowAccessor) Ref() lineRef {
	return lineRef{lineKindRow, acc.index}
}

type colAccessor struct {
	index int
	board *Board
}

func (acc colAccessor) Get() []Cell {
	cells := make([]Cell, acc.board.GetRows())
	for i := range *acc.board {
		cells[i] = (*acc.board)[i][acc.index]
	}
	return cells
}

func (acc colAccessor) Set(cells []Cell) {
	for i := range cells {
		(*acc.board)[i][acc.index] = cells[i]
	}
}

func (acc colAccessor) Ref() lineRef {
	return lineRef{lineKindColumn, acc.index}
}
