package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

func (r MinimumSpacingRule) Name() string {
	return "MinimumSpacingRule"
}

func (r MinimumSpacingRule) Deduce(line line.Line) []game.Cell {
	cells := slices.Clone(line.Cells)

	segs := SplitByWhite(cells)
	if len(segs) != 1 {
		return nil
	}

	seg := segs[0]
	var sum int
	for _, h := range line.Hints {
		sum += h
	}
	if sum+(len(line.Hints)-1) != len(seg) {
		return nil
	}

	var last int
	for i, hint := range line.Hints {
		for range hint {
			seg[last] = game.CellBlack
			last++
		}
		if i != len(line.Hints)-1 {
			seg[last] = game.CellWhite
			last++
		}
	}
	return cells
}
