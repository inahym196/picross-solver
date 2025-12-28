package picrosssolver

import "strings"

type Cell uint8

const (
	CellUnknown Cell = iota
	CellWhite
	CellBlack
)

func (c Cell) String() string {
	switch c {
	case CellUnknown:
		return "CellUnknown"
	case CellBlack:
		return "CellBlack"
	case CellWhite:
		return "CellWhite"
	default:
		panic("invalid cell")
	}
}

func filledCells(length int, c Cell) []Cell {
	line := make([]Cell, length)
	for i := range line {
		line[i] = c
	}
	return line
}

type Board [][]Cell

func newBoard(h, w int) Board {
	board := make(Board, h)
	for i := range h {
		board[i] = make([]Cell, w)
	}
	return board
}

type LineType uint8

const (
	LineRow LineType = iota
	LineColumn
)

func (b Board) Print() []string {
	var ss []string
	for i := range b {
		var s strings.Builder
		for j := range b[i] {
			switch b[i][j] {
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

func DeepCopyBoard(src Board) Board {
	dst := make(Board, len(src))
	for i := range src {
		dst[i] = make([]Cell, len(src[i]))
		copy(dst[i], src[i])
	}
	return dst
}

type Game struct {
	board    Board
	rowHints [][]int
	colHints [][]int
}

func NewGame(rowHints, colHints [][]int) *Game {
	b := newBoard(len(rowHints), len(colHints))
	return &Game{b, rowHints, colHints}
}

type Line struct {
	Cells     []Cell
	Hints     []int
	WriteBack func([]Cell)
}

func (l Line) IsAllCells(c Cell) bool {
	for _, cell := range l.Cells {
		if cell != c {
			return false
		}
	}
	return true
}

type Rule interface {
	Apply(line Line)
}

type FillCompleteRule struct{}

func (r *FillCompleteRule) Apply(line Line) {
	if len(line.Hints) != 1 {
		return
	}

	hint := line.Hints[0]
	switch {
	case hint == len(line.Cells):
		if !line.IsAllCells(CellBlack) {
			updated := filledCells(len(line.Cells), CellBlack)
			line.WriteBack(updated)
		}
	case hint == 0:
		if !line.IsAllCells(CellWhite) {
			updated := filledCells(len(line.Cells), CellWhite)
			line.WriteBack(updated)
		}
	}
}

type Solver struct {
	rules []Rule
}

func NewSolver() Solver {
	rules := []Rule{&FillCompleteRule{}}
	return Solver{rules}
}

func Lines(board Board, rowHints, colHints [][]int) []Line {
	var lines []Line

	for i := range board {
		lines = append(lines, Line{
			Cells: board[i],
			Hints: rowHints[i],
			WriteBack: func(updated []Cell) {
				copy(board[i], updated)
			},
		})
	}
	for i := range len(board[0]) {
		col := make([]Cell, len(board))
		lines = append(lines, Line{
			Cells: col,
			Hints: colHints[i],
			WriteBack: func(updated []Cell) {
				for r := range board {
					board[r][i] = updated[r]
				}
			},
		})
	}
	return lines
}

func (s Solver) ApplyOnce(game Game) Board {
	board := DeepCopyBoard(game.board)
	lines := Lines(board, game.rowHints, game.colHints)
	for _, rule := range s.rules {
		for _, line := range lines {
			rule.Apply(line)
		}
	}
	return board
}
