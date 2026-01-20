package game

import (
	"errors"
	"fmt"
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

type board struct {
	cells  [][]Cell
	width  int
	height int
}

func newBoard(width, height int) *board {
	cells := make([][]Cell, height)
	for i := range height {
		cells[i] = make([]Cell, width)
	}
	return &board{cells, width, height}
}

func (b *board) InBounds(row, col int) bool {
	return 0 <= row && row < len(b.cells) && 0 <= col && col < len(b.cells[0])
}

func (b *board) SetCell(row, col int, cell Cell) {
	b.cells[row][col] = cell
}

func (b *board) Cell(row, col int) (Cell, error) {
	if !b.InBounds(row, col) {
		return CellUndetermined, fmt.Errorf("out of range")
	}
	return b.cells[row][col], nil
}

func (b *board) Print() []string {
	ss := make([]string, 0, len(b.cells))
	for i := range b.cells {
		var s strings.Builder
		s.Grow(len(b.cells[i]))
		for _, c := range b.cells[i] {
			switch c {
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

type LineRef struct {
	kind  LineKind
	index int
}

func (ref LineRef) String() string {
	return fmt.Sprintf("%s[%d]", ref.kind, ref.index)
}

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

type LineView struct {
	board   *board
	lineRef LineRef
	hints   []int
}

type Game struct {
	board    *board
	RowHints [][]int
	ColHints [][]int
}

func NewGame(rowHints, colHints [][]int) (*Game, error) {
	if len(rowHints) == 0 || len(colHints) == 0 {
		return nil, errors.New("rowHints,colHintsは1より大きい必要がある")
	}
	// TODO: hintsの最小配置がlen(cell)より小さい必要がある
	width := len(colHints)
	height := len(rowHints)

	b := newBoard(width, height)
	return &Game{b, rowHints, colHints}, nil
}

func (g *Game) SetCell(row, col int, cell Cell) error {
	if !g.board.InBounds(row, col) {
		return fmt.Errorf("out of range")
	}
	g.board.SetCell(row, col, cell)
	return nil
}

func (g *Game) Print() []string {
	return g.board.Print()
}

func (g *Game) Row(i int) LineView {
	return LineView{g.board, LineRef{LineKindRow, i}, g.RowHints[i]}
}

func (g *Game) Col(i int) LineView {
	return LineView{g.board, LineRef{LineKindColumn, i}, g.ColHints[i]}
}
