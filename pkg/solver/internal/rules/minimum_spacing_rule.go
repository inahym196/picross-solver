package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

func (r MinimumSpacingRule) Name() string {
	return "MinimumSpacingRule"
}

func (r MinimumSpacingRule) Deduce(line game.Line) []game.Cell {
	cells := slices.Clone(line.Cells)
	trim := trimWhite(cells)

	var sum int
	for _, h := range line.Hints {
		sum += h
	}
	if sum+(len(line.Hints)-1) != len(trim) {
		return nil
	}

	last := 0
	for i, hint := range line.Hints {
		for range hint {
			trim[last] = game.CellBlack
			last++
		}
		if i != len(line.Hints)-1 {
			trim[last] = game.CellWhite
			last++
		}
	}
	return cells
}
