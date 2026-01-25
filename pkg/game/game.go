package game

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Cell uint8

const (
	CellUndetermined Cell = iota
	CellWhite
	CellBlack
)

func (c Cell) String() string {
	switch c {
	case CellUndetermined:
		return "U"
	case CellBlack:
		return "B"
	case CellWhite:
		return "W"
	default:
		panic("invalid cell")
	}
}

// Entity
type Board struct {
	cells  [][]Cell
	width  int
	height int
}

func NewBoard(width, height int) *Board {
	cells := make([][]Cell, height)
	for i := range height {
		cells[i] = make([]Cell, width)
	}
	return &Board{cells, width, height}
}

func (b *Board) Width() int  { return len(b.cells) }
func (b *Board) Height() int { return len(b.cells[0]) }

func (b *Board) Row(i int) []Cell { return slices.Clone(b.cells[i]) }
func (b *Board) Col(i int) []Cell {
	cells := make([]Cell, b.height)
	for row := range b.height {
		cells[i] = b.cells[row][i]
	}
	return cells
}

func (b *Board) Cells() [][]Cell {
	h := len(b.cells)
	cells := make([][]Cell, h)
	for i := range h {
		cells[i] = slices.Clone(b.cells[i])
	}
	return cells
}

func (b *Board) Mark(row, col int, cell Cell) error {
	if !b.inBounds(row, col) {
		return fmt.Errorf("out of range")
	}
	b.cells[row][col] = cell
	return nil
}

func (b *Board) inBounds(row, col int) bool {
	return 0 <= row && row < b.height && 0 <= col && col <= b.width
}

func (b *Board) Print() []string {
	var ss []string
	for i := range b.cells {
		var s strings.Builder
		for j := range b.cells[i] {
			switch b.cells[i][j] {
			case CellBlack:
				s.WriteString("#")
			case CellWhite:
				s.WriteString("_")
			case CellUndetermined:
				s.WriteString("?")
			}
		}
		ss = append(ss, s.String())
	}
	return ss
}

type LineKind uint8

const (
	LineKindRow LineKind = iota
	LineKindColumn
)

func (kind LineKind) String() string {
	switch kind {
	case LineKindRow:
		return "Row"
	case LineKindColumn:
		return "Col"
	default:
		panic("invalid lineKind")
	}
}

type LineRef struct {
	Kind  LineKind
	Index int
}

func (ref LineRef) String() string {
	return fmt.Sprintf("%s[%d]", ref.Kind, ref.Index)
}

type Game struct {
	board    *Board
	RowHints [][]int
	ColHints [][]int
}

func NewGame(rowHints, colHints [][]int) (Game, error) {
	if len(rowHints) == 0 || len(colHints) == 0 {
		return Game{}, errors.New("rowHints,colHintsは1より大きい必要がある")
	}
	// TODO: hintsの最小配置がlen(cell)より小さい必要がある
	width := len(colHints)
	height := len(rowHints)

	b := NewBoard(height, width)
	return Game{b, rowHints, colHints}, nil
}

func (g Game) Board() *Board {
	return g.board
}

func (g Game) Mark(row, col int, cell Cell) error {
	return g.board.Mark(row, col, cell)
}
