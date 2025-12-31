package picrosssolver

import (
	"errors"
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

func NewGame(rowHints, colHints [][]int) (*Game, error) {
	if len(rowHints) == 0 || len(colHints) == 0 {
		return nil, errors.New("rowHints,colHintsは1より大きい必要がある")
	}
	// TODO: hintsの最小配置がlen(cell)より小さい必要がある
	width := len(colHints)
	height := len(rowHints)

	b := newBoard(height, width)
	return &Game{b, rowHints, colHints}, nil
}
