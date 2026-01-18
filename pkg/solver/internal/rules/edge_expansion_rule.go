package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) Name() string {
	return "EdgeExpansionRule"
}

func (r EdgeExpansionRule) applyLeft(cells []game.Cell, hint int) (changed bool) {
	seg := SplitByWhite(cells)[0]
	if seg[0] != game.CellBlack || len(seg) < hint {
		return false
	}
	for i := range hint {
		seg[i] = game.CellBlack
		changed = true
	}
	if len(seg) > hint {
		seg[hint] = game.CellWhite
	}
	return changed
}

func (r EdgeExpansionRule) Deduce(line line.Line) []game.Cell {
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
