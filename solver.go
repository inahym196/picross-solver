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

type ExtractMatchRule struct{}

func (r *ExtractMatchRule) Apply(line Line) {
	if len(line.Hints) != 1 {
		return
	}

	hint := line.Hints[0]
	if hint == len(line.Cells) && !line.IsAllCells(CellBlack) {
		updated := filledCells(len(line.Cells), CellBlack)
		line.WriteBack(updated)
	}
}

type ZeroHintRule struct{}

func (r *ZeroHintRule) Apply(line Line) {
	if len(line.Hints) != 1 {
		return
	}

	hint := line.Hints[0]
	if hint == 0 && !line.IsAllCells(CellWhite) {
		updated := filledCells(len(line.Cells), CellWhite)
		line.WriteBack(updated)
	}
}

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

// ヒントブロックを左詰め／右詰めしたときに必ず重なる部分を黒確定
type OverlapFillRule struct{}

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpantionRule struct{}

// 既に黒が hint 長に達しているブロックの前後を白確定
type BlockSatisfiedRule struct{}

// 白確定セルで line を分割し、それぞれにヒントを再配分
type SegmentSplitRule struct{}

// ヒントが収まらない区間を白確定
type PruneImpossibleSegmentRule struct{}

// 長さ < 最小 hint の区間はすべて白
type TooSmallSegmentRule struct{}

// すべての hint を満たした後の残りは白
type FillRemainingWhiteRule struct{}

// 仮に黒／白を置き、矛盾が出たら逆を確定
type HypothesisRule struct{}

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
	lines := s.ExtractLines(board, game.rowHints, game.colHints)
	for _, rule := range s.rules {
		for _, line := range lines {
			rule.Apply(line)
		}
	}
	return board
}
