package picrosssolver

import "strings"

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

func Solve(rowHints, colHints [][]int) Board {
	w, h := len(colHints), len(rowHints)
	b := newBoard(h, w)
	for i := range h {
		if len(rowHints[i]) == 1 {
			switch rowHints[i][0] {
			case h:
				b.paintLine(LineRow, i, CellBlack)
			case 0:
				b.paintLine(LineRow, i, CellWhite)
			}
		}
	}
	for i := range w {
		if len(colHints[i]) == 1 {
			switch colHints[i][0] {
			case w:
				b.paintLine(LineColumn, i, CellBlack)
			case 0:
				b.paintLine(LineColumn, i, CellWhite)
			}
		}
	}
	return b
}
