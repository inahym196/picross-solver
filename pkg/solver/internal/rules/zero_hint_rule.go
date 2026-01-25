package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

type ZeroHintRule struct{}

func (e ZeroHintRule) Name() string {
	return "ZeroHintRule"
}

func (r ZeroHintRule) Deduce(line game.Line) []game.Cell {
	cells := slices.Clone(line.Cells)
	if len(line.Hints) != 1 || line.Hints[0] != 0 {
		return nil
	}
	for i := range cells {
		cells[i] = game.CellWhite
	}
	return cells
}
