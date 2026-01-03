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

type lineAccessor interface {
	Get() []Cell
	Set(cells []Cell)
}

type rowAccessor struct {
	index int
	board *Board
}

func (acc rowAccessor) Get() []Cell {
	return slices.Clone((*acc.board)[acc.index])
}

func (acc rowAccessor) Set(cells []Cell) {
	copy((*acc.board)[acc.index], cells)
}

type colAccessor struct {
	index int
	board *Board
}

func (acc colAccessor) Get() []Cell {
	cells := make([]Cell, acc.board.GetRows())
	for i := range *acc.board {
		cells[i] = (*acc.board)[i][acc.index]
	}
	return cells
}

func (acc colAccessor) Set(cells []Cell) {
	for i := range cells {
		(*acc.board)[i][acc.index] = cells[i]
	}
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
		cells := acc.Get()
		if slices.Index(cells, CellUndetermined) == -1 {
			return
		}
		hc := NewHintedCells(slices.Clone(cells), slices.Clone(hints))
		updated := rule.Deduce(hc)
		if updated != nil && !reflect.DeepEqual(cells, updated) {
			acc.Set(updated)
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
