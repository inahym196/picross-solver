package picrosssolver

import (
	"fmt"
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

type applyLog struct {
	ruleName string
	lineRef  lineRef
	before   []Cell
	after    []Cell
}

func (log applyLog) String() string {
	return fmt.Sprintf("%s %s %v -> %v", log.ruleName, log.lineRef, log.before, log.after)
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

func (s Solver) ApplyLine(acc lineAccessor, hints []int) []applyLog {
	var logs []applyLog

	// TODO: lineごとにrulesを適用し、最後にApplyすればApply頻度を下げられる
	for _, rule := range s.rules {
		before := acc.get()
		if slices.Index(before, CellUndetermined) == -1 {
			return logs
		}
		hc := NewHintedCells(slices.Clone(before), slices.Clone(hints))
		updated := rule.Deduce(hc)
		if updated != nil && !reflect.DeepEqual(before, updated) {
			acc.set(updated)
			logs = append(logs, applyLog{
				ruleName: rule.Name(),
				lineRef:  acc.ref(),
				before:   before,
				after:    updated,
			})
		}
	}
	return logs
}

func (s Solver) ApplyOnce(game Game) (Board, []applyLog) {
	var logs []applyLog

	board := slices.Clone(game.board)
	for i := range game.rowHints {
		lineLogs := s.ApplyLine(rowAccessor{i, &board}, game.rowHints[i])
		logs = append(logs, lineLogs...)
	}
	for i := range game.colHints {
		lineLogs := s.ApplyLine(colAccessor{i, &board}, game.colHints[i])
		logs = append(logs, lineLogs...)
	}
	return board, logs
}

func (s Solver) checkComplete(board Board) bool {
	for row := range board {
		if slices.Index(board[row], CellUndetermined) != -1 {
			return false
		}
	}
	return true
}

func (s Solver) ApplyMany(game Game) (Board, int, []applyLog) {
	var logs []applyLog
	board := DeepCopyBoard(game.board)
	n := 0
	for !s.checkComplete(board) {
		n++
		deduced, lineLogs := s.ApplyOnce(game)
		if reflect.DeepEqual(board, deduced) {
			return board, n, logs
		}
		board = deduced
		logs = append(logs, lineLogs...)
	}
	return board, n, logs
}
