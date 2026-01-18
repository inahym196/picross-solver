package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

// 端が未確定なら黒をヒント分拡張する
type OverlapExpansionRule struct{}

func (r OverlapExpansionRule) Name() string {
	return "OverlapExpansionRule"
}

func (r OverlapExpansionRule) applyLeft(cells []game.Cell, hint int) (changed bool) {
	seg := SplitByWhite(cells)[0]
	firstBlackIndex := slices.Index(seg, game.CellBlack)
	if firstBlackIndex == -1 || firstBlackIndex >= hint {
		return false
	}

	for i := firstBlackIndex + 1; i < hint; i++ {
		// TODO: バグの可能性あり
		seg[i] = game.CellBlack
		changed = true
	}
	return changed
}

func (r OverlapExpansionRule) Deduce(line line.Line) []game.Cell {
	cells := slices.Clone(line.Cells)

	firstHint := line.Hints[0]
	changed1 := r.applyLeft(cells, firstHint)

	slices.Reverse(cells)
	lastHint := line.Hints[len(line.Hints)-1]
	changed2 := r.applyLeft(cells, lastHint)

	if !changed1 && !changed2 {
		return nil
	}
	slices.Reverse(cells)
	return cells
}
