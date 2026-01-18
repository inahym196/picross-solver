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

func (r OverlapExpansionRule) expand(cells []game.Cell, hint int, indexer func(i, n int) int) bool {
	n := len(cells)

	firstBlackIndex := -1
	for i := range n {
		if cells[indexer(i, n)] == game.CellBlack {
			firstBlackIndex = i
			break
		}
	}
	if firstBlackIndex == -1 || firstBlackIndex >= hint {
		return false
	}

	for i := firstBlackIndex + 1; i < hint; i++ {
		cells[indexer(i, n)] = game.CellBlack
	}
	return true
}

func (r OverlapExpansionRule) Deduce(line line.Line) []game.Cell {
	cells := slices.Clone(line.Cells)
	trim := trimWhite(cells)

	leftIndexer := func(i, _ int) int { return i }
	expand1 := r.expand(trim, line.Hints[0], leftIndexer)

	rightIndexer := func(i, n int) int { return n - 1 - i }
	expand2 := r.expand(cells, line.Hints[len(line.Hints)-1], rightIndexer)

	if !expand1 && !expand2 {
		return nil
	}
	return cells
}
