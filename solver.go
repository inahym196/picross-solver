package picrosssolver

import (
	"reflect"
	"slices"
)

type Solver struct {
	deducer deducer
}

func NewSolver() Solver {
	return Solver{newDeducer()}
}

func (s Solver) ApplyLine(acc lineAccessor, hints []int) (changed bool, deductions []deduction) {
	// TODO: lineごとにrulesを適用し、最後にApplyすればApply頻度を下げられる
	for _, rule := range s.deducer.rules {
		before := acc.Cells()
		if slices.Index(before, CellUndetermined) == -1 {
			return changed, deductions
		}
		line := lineView{slices.Clone(before), slices.Clone(hints)}
		updated := rule.Deduce(line)
		if updated != nil && !reflect.DeepEqual(before, updated) {
			changed = true
			acc.Update(updated)
			deductions = append(deductions, deduction{
				ruleName: rule.Name(),
				hints:    hints,
				lineRef:  acc.Ref(),
				before:   before,
				after:    updated,
			})
		}
	}
	return changed, deductions
}

func (s Solver) ApplyOnce(game *Game) (board Board, changed bool, deductions []deduction) {

	board = slices.Clone(game.board)
	for i := range game.rowHints {
		acc := lineAccessor{&game.board, lineRef{lineKindRow, i}}
		changedLine, lineDeds := s.ApplyLine(acc, game.rowHints[i])
		if changedLine {
			deductions = append(deductions, lineDeds...)
			changed = true
		}
	}
	for i := range game.colHints {
		acc := lineAccessor{&game.board, lineRef{lineKindColumn, i}}
		changedLine, lineDeds := s.ApplyLine(acc, game.colHints[i])
		if changedLine {
			deductions = append(deductions, lineDeds...)
			changed = true
		}
	}
	return board, changed, deductions
}

func (s Solver) checkComplete(board Board) bool {
	for row := range board {
		if slices.Index(board[row], CellUndetermined) != -1 {
			return false
		}
	}
	return true
}

func (s Solver) ApplyMany(game *Game) (Board, int, []deduction) {
	var deds []deduction
	board := DeepCopyBoard(game.board)
	n := 0
	for !s.checkComplete(board) {
		n++
		deduced, changedLine, lineDeds := s.ApplyOnce(game)
		if !changedLine {
			return board, n, deds
		}
		board = deduced
		deds = append(deds, lineDeds...)
	}
	return board, n, deds
}
