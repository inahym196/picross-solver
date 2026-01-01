package picrosssolver

import (
	"reflect"
)

type LineKind uint8

const (
	LineKindRow LineKind = iota
	LineKindColumn
)

func (lk LineKind) String() string {
	switch lk {
	case LineKindRow:
		return "Row"
	case LineKindColumn:
		return "Column"
	default:
		panic("invalid LineKind")
	}

}

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
}

type lineAccessor struct {
	Get func() []Cell
	Set func(cells []Cell)
}

func rowAccessor(board Board, row int) lineAccessor {
	return lineAccessor{
		Get: func() []Cell {
			cells := make([]Cell, board.GetColumns())
			copy(cells, board[row])
			return cells
		},
		Set: func(cells []Cell) { copy(board[row], cells) },
	}
}

func colAccessor(board Board, col int) lineAccessor {
	return lineAccessor{
		Get: func() []Cell {
			cells := make([]Cell, board.GetRows())
			for i := range board {
				cells[i] = board[i][col]
			}
			return cells
		},
		Set: func(cells []Cell) {
			for i := range cells {
				board[i][col] = cells[i]
			}
		},
	}
}

type Solver struct {
	rules []Rule
}

func NewSolver() Solver {
	rules := []Rule{
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

func (s Solver) ApplyLine(acc lineAccessor, hints []int) {
	// TODO: lineごとにrulesを適用し、最後にApplyすればApply頻度を下げられる
	for _, rule := range s.rules {
		hc := NewHintedCells(acc.Get(), hints)
		updated := rule.Deduce(hc)
		if updated != nil && !reflect.DeepEqual(acc.Get(), updated) {
			acc.Set(updated)
		}
	}
}

func (s Solver) ApplyOnce(game Game) Board {
	board := DeepCopyBoard(game.board)
	for i := range game.rowHints {
		acc := rowAccessor(board, i)
		s.ApplyLine(acc, game.rowHints[i])
	}
	for i := range game.colHints {
		acc := colAccessor(board, i)
		s.ApplyLine(acc, game.colHints[i])
	}
	return board
}
