package picrosssolver

import "reflect"

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

func DeepCopyHintedCells(hc HintedCells) HintedCells {
	cells := make([]Cell, len(hc.Cells))
	copy(cells, hc.Cells)
	hints := make([]int, len(hc.Hints))
	copy(hints, hc.Hints)
	return NewHintedCells(cells, hints)
}

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
		ExtractMatchRule{},
		ZeroHintRule{},
		MinimumSpacingRule{},
		//OverlapFillRule{},
		//OverlapExpansionRule{},
		//EdgeExpansionRule{},
		//BlockSatisfiedRule{},
		//PruneImpossibleSegmentRule{},
		//FillRemainingWhiteRule{},
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
			hc := DeepCopyHintedCells(line.Cells)
			updated := rule.Deduce(hc)
			if updated != nil && !reflect.DeepEqual(line.Cells, updated) {
				s.ApplyLine(board, line, updated)
			}
		}
	}
	return board
}
