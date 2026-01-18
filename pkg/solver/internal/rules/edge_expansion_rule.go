package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

func trimWhite(cells []game.Cell) []game.Cell {
	start := 0
	for start < len(cells) && cells[start] == game.CellWhite {
		start++
	}
	end := len(cells) - 1
	for end > start && cells[end] == game.CellWhite {
		end--
	}
	if start >= end {
		return nil
	}
	return cells[start : end+1]
}

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) Name() string {
	return "EdgeExpansionRule"
}

func (r EdgeExpansionRule) expand(cells []game.Cell, hint int, idx func(i, n int) int) bool {
	n := len(cells)
	if n < hint || cells[idx(0, n)] != game.CellBlack {
		return false
	}

	for i := range hint {
		cells[idx(i, n)] = game.CellBlack
	}
	if hint < n {
		cells[idx(hint, n)] = game.CellWhite
	}
	return true
}

func (r EdgeExpansionRule) Deduce(line line.Line) []game.Cell {
	cells := slices.Clone(line.Cells)

	trim := trimWhite(cells)

	leftIndex := func(i, _ int) int { return i }
	expand1 := r.expand(trim, line.Hints[0], leftIndex)

	rightIndex := func(i, n int) int { return n - 1 - i }
	expand2 := r.expand(trim, line.Hints[len(line.Hints)-1], rightIndex)

	if !expand1 && !expand2 {
		return nil
	}
	return cells
}
