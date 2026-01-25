package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

// すべての hint を満たした後の残りは白
type FillRemainingWhiteRule struct{}

func (r FillRemainingWhiteRule) Name() string {
	return "FillRemainingWhiteRule"
}

func (r FillRemainingWhiteRule) Deduce(line game.Line) []game.Cell {
	cells := slices.Clone(line.Cells)

	sumHints := 0
	for _, h := range line.Hints {
		sumHints += h
	}

	blackCount := 0
	for _, c := range cells {
		if c == game.CellBlack {
			blackCount++
		}
	}

	if blackCount != sumHints {
		return nil
	}

	changed := false
	for i, c := range cells {
		if c == game.CellUndetermined {
			cells[i] = game.CellWhite
			changed = true
		}
	}

	if !changed {
		return nil
	}
	return cells
}
