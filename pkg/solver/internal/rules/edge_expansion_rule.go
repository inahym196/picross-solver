package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) Name() string {
	return "EdgeExpansionRule"
}

func (r EdgeExpansionRule) expand(cells []game.Cell, hint int, indexer func(i, n int) int) bool {
	n := len(cells)

	if n < hint || cells[indexer(0, n)] != game.CellBlack {
		return false
	}

	for i := range hint {
		cells[indexer(i, n)] = game.CellBlack
	}
	if hint < n {
		cells[indexer(hint, n)] = game.CellWhite
	}
	return true
}

func (r EdgeExpansionRule) Deduce(line game.Line) []game.Cell {
	cells := slices.Clone(line.Cells)
	trim := trimWhite(cells)

	leftIndexer := func(i, _ int) int { return i }
	expand1 := r.expand(trim, line.Hints[0], leftIndexer)

	rightIndexer := func(i, n int) int { return n - 1 - i }
	expand2 := r.expand(trim, line.Hints[len(line.Hints)-1], rightIndexer)

	if !expand1 && !expand2 {
		return nil
	}
	return cells
}
