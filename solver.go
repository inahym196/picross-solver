package picrosssolver

import (
	"reflect"
	"slices"
)

type HintedCells struct {
	Cells []Cell
	Hints []int
}

func NewHintedCells(cells []Cell, hints []int) HintedCells {
	return HintedCells{cells, hints}
}

type Solver struct {
	rules []Rule
}

func NewSolver() Solver {
	rules := []Rule{
		ZeroHintRule{},
		MinimumSpacingRule{},
		OverlapFillRule{},
		OverlapExpansionRule{},
		EdgeExpansionRule{},
		BlockSatisfiedRule{},
		PruneImpossibleSegmentRule{},
		FillRemainingWhiteRule{},
	}
	return Solver{rules}
}

func (s Solver) ApplyLine(acc lineAccessor, hints []int) {
	// TODO: lineごとにrulesを適用し、最後にApplyすればApply頻度を下げられる
	for _, rule := range s.rules {
		cells := acc.get()
		if slices.Index(cells, CellUndetermined) == -1 {
			return
		}
		hc := NewHintedCells(slices.Clone(cells), slices.Clone(hints))
		updated := rule.Deduce(hc)
		if updated != nil && !reflect.DeepEqual(cells, updated) {
			acc.set(updated)
		}
	}
}

func (s Solver) ApplyOnce(game Game) Board {
	board := slices.Clone(game.board)
	for i := range game.rowHints {
		acc := rowAccessor{i, &board}
		s.ApplyLine(acc, game.rowHints[i])
	}
	for i := range game.colHints {
		acc := colAccessor{i, &board}
		s.ApplyLine(acc, game.colHints[i])
	}
	return board
}

func (s Solver) checkComplete(board Board) bool {
	for row := range board {
		if slices.Index(board[row], CellUndetermined) != -1 {
			return false
		}
	}
	return true
}

func (s Solver) ApplyMany(game Game) (Board, int) {
	board := DeepCopyBoard(game.board)
	n := 0
	for !s.checkComplete(board) {
		n++
		deduced := s.ApplyOnce(game)
		if reflect.DeepEqual(board, deduced) {
			return board, n
		}
		board = deduced
	}
	return board, n
}
