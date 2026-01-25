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

type Board struct {
	cells [][]Cell
}

func newBoard(height, width int) Board {
	cells := make([][]Cell, height)
	for i := range height {
		cells[i] = make([]Cell, width)
	}
	return Board{cells}
}

func (b Board) InBounds(row, col int) bool {
	return 0 <= row && row < len(b.cells) && 0 <= col && col <= len(b.cells[0])
}

func (b Board) SetCell(row, col int, cell Cell) {
	b.cells[row][col] = cell
}

func (b Board) GetRows() int {
	return len(b.cells)
}

func (b Board) GetColumns() int {
	return len(b.cells[0])
}

func (b Board) Cells() [][]Cell {
	h := len(b.cells)
	cells := make([][]Cell, h)
	for i := range h {
		cells[i] = slices.Clone(b.cells[i])
	}
	return cells
}

func (b Board) Print() []string {
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

type Game struct {
	board    Board
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

	b := newBoard(height, width)
	return Game{b, rowHints, colHints}, nil
}

func (g Game) Board() Board {
	return g.board
}

func (g Game) SetCell(row, col int, cell Cell) error {
	if !g.board.InBounds(row, col) {
		return fmt.Errorf("out of range")
	}
	g.board.SetCell(row, col, cell)
	return nil
}
