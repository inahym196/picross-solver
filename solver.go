package picrosssolver

import (
	"fmt"
	"reflect"
	"slices"
)

type lineView struct {
	Cells []Cell
	Hints []int
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
		before := acc.Cells()
		if slices.Index(before, CellUndetermined) == -1 {
			return changed, logs
		}
		line := lineView{slices.Clone(before), slices.Clone(hints)}
		updated := rule.Deduce(line)
		if updated != nil && !reflect.DeepEqual(before, updated) {
			changed = true
			acc.Update(updated)
			logs = append(logs, applyLog{
				ruleName: rule.Name(),
				hints:    hints,
				lineRef:  acc.Ref(),
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
		acc := lineAccessor{&game.board, lineRef{lineKindRow, i}}
		changedLine, lineLogs := s.ApplyLine(acc, game.rowHints[i])
		if changedLine {
			logs = append(logs, lineLogs...)
			changed = true
		}
	}
	for i := range game.colHints {
		acc := lineAccessor{&game.board, lineRef{lineKindColumn, i}}
		changedLine, lineLogs := s.ApplyLine(acc, game.colHints[i])
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
