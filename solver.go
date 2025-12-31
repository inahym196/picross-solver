package picrosssolver

import (
	"reflect"
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
		return "CellUndetermined"
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

func newBoard(height, width int) Board {
	board := make(Board, height)
	for i := range height {
		board[i] = make([]Cell, width)
	}
	return board
}

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
			case CellUndetermined:
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
	// TODO: Hintsはそれぞれ1より大きいという制約を入れる
	// 仮置きでpanicさせておく
	if len(rowHints) < 1 || len(colHints) < 1 {
		panic("rowHints,colHintsは1より大きい必要がある")
	}
	b := newBoard(len(rowHints), len(colHints))
	return &Game{b, rowHints, colHints}
}

type LineKind uint8

const (
	LineKindRow LineKind = iota
	LineKindColumn
)

type HintedCells struct {
	Cells []Cell
	Hints []int
}

func NewHintedCells(cells []Cell, hints []int) HintedCells { return HintedCells{cells, hints} }

type Line struct {
	Kind  LineKind
	Index int
	Cells HintedCells
}

type Solver struct {
	rules []Rule
}

func NewSolver() Solver {
	rules := []Rule{
		&ExtractMatchRule{},
		&ZeroHintRule{},
	}
	return Solver{rules}
}

func (s Solver) ExtractLines(board Board, rowHints, colHints [][]int) []Line {
	var lines []Line

	for i := range board {
		lines = append(lines, Line{
			Kind:  LineKindRow,
			Index: i,
			Cells: NewHintedCells(board[i], rowHints[i]),
		})
	}
	for i := range board[0] {
		col := make([]Cell, len(board[0]))
		for row := range board {
			col[row] = board[row][i]
		}
		lines = append(lines, Line{
			Kind:  LineKindColumn,
			Index: i,
			Cells: NewHintedCells(col, colHints[i]),
		})
	}
	return lines
}

func (s Solver) ApplyLine(board Board, line Line, cells []Cell) {
	switch line.Kind {
	case LineKindRow:
		copy(board[line.Index], cells)
	case LineKindColumn:
		for row := range board {
			board[row][line.Index] = cells[row]
		}
	default:
		panic("invalid LineKind")
	}
}

func (s Solver) ApplyOnce(game Game) Board {
	board := DeepCopyBoard(game.board)
	lines := s.ExtractLines(board, game.rowHints, game.colHints)
	for _, rule := range s.rules {
		for _, line := range lines {
			// TODO: lineごとにrulesを適用し、最後にApplyすればApply頻度を下げられる
			// TODO: rule内でline.Cellを破壊されないようにcopyを渡した方がいい
			updated := rule.Deduce(line.Cells)
			if updated != nil && !reflect.DeepEqual(line.Cells, updated) {
				s.ApplyLine(board, line, updated)
			}
		}
	}
	return board
}
