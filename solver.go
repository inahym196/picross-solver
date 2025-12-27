package picrosssolver

import "strings"

type Hints struct {
	row [][]int
	col [][]int
}

func NewHints(row, col [][]int) Hints {
	h := Hints{row, col}
	return h
}

type Cell uint8

const (
	CellUnknown Cell = iota
	CellWhite
	CellBlack
)

type Board struct {
	cells [][]Cell
}

func newBoard(h, w int) Board {
	cells := make([][]Cell, h)
	for i := range h {
		cells[i] = make([]Cell, w)
	}
	return Board{cells}
}

type LineType uint8

const (
	LineRow LineType = iota
	LineColumn
)

func (b *Board) paintLine(typ LineType, num int, cell Cell) {
	switch typ {
	case LineRow:
		for i := range b.cells[num] {
			b.cells[num][i] = cell
		}
	case LineColumn:
		for i := range len(b.cells) {
			b.cells[i][num] = cell
		}
	}
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
			case CellUnknown:
				s.WriteString("?")
			}
		}
		ss = append(ss, s.String())
	}
	return ss
}

func Solve(hints Hints) Board {
	w, h := len(hints.col), len(hints.row)
	b := newBoard(h, w)
	for i := range h {
		if len(hints.row[i]) == 1 {
			switch hints.row[i][0] {
			case h:
				b.paintLine(LineRow, i, CellBlack)
			case 0:
				b.paintLine(LineRow, i, CellWhite)
			}
		}
	}
	for i := range w {
		if len(hints.col[i]) == 1 {
			switch hints.col[i][0] {
			case w:
				b.paintLine(LineColumn, i, CellBlack)
			case 0:
				b.paintLine(LineColumn, i, CellWhite)
			}
		}
	}
	return b
}
