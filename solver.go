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
	hints    []int
	lineRef  lineRef
	before   []Cell
	after    []Cell
}

func (log applyLog) String() string {
	return fmt.Sprintf("%s %s %v %v -> %v", log.ruleName, log.lineRef, log.hints, log.before, log.after)
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

func (s Solver) ApplyLine(acc lineAccessor, hints []int) (changed bool, logs []applyLog) {
	// TODO: lineごとにrulesを適用し、最後にApplyすればApply頻度を下げられる
	for _, rule := range s.rules {
		before := acc.get()
		if slices.Index(before, CellUndetermined) == -1 {
			return changed, logs
		}
		hc := NewHintedCells(slices.Clone(before), slices.Clone(hints))
		updated := rule.Deduce(hc)
		if updated != nil && !reflect.DeepEqual(before, updated) {
			changed = true
			acc.set(updated)
			logs = append(logs, applyLog{
				ruleName: rule.Name(),
				hints:    hints,
				lineRef:  acc.ref(),
				before:   before,
				after:    updated,
			})
		}
	}
	return changed, logs
}

func (s Solver) ApplyOnce(game Game) (board Board, changed bool, logs []applyLog) {

	board = slices.Clone(game.board)
	for i := range game.rowHints {
		changedLine, lineLogs := s.ApplyLine(rowAccessor{i, &board}, game.rowHints[i])
		if changedLine {
			logs = append(logs, lineLogs...)
			changed = true
		}
	}
	for i := range game.colHints {
		changedLine, lineLogs := s.ApplyLine(colAccessor{i, &board}, game.colHints[i])
		if changedLine {
			logs = append(logs, lineLogs...)
			changed = true
		}
	}
	return board, changed, logs
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
		deduced, changedLine, lineLogs := s.ApplyOnce(game)
		if !changedLine {
			return board, n, logs
		}
		board = deduced
		logs = append(logs, lineLogs...)
	}
	return board, n, logs
}
