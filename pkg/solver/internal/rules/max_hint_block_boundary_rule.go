package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

// 最大 hint 長に達しているブロックの前後を白確定
type MaxHintBlockBoundaryRule struct{}

func (r MaxHintBlockBoundaryRule) Name() string {
	return "BlockSatisfiedRule"
}

func (r MaxHintBlockBoundaryRule) Deduce(line game.Line) []game.Cell {
	cells := slices.Clone(line.Cells)
	hint := slices.Max(line.Hints)
	if hint == 0 {
		return nil
	}

	blocks := findBlocksN(cells, hint)
	if len(blocks) == 0 {
		return nil
	}

	changed := false
	for _, block := range blocks {
		prevStart := block.start - 1
		if prevStart >= 0 && cells[prevStart] == game.CellUndetermined {
			cells[prevStart] = game.CellWhite
			changed = true
		}
		afterEnd := block.start + block.length
		if afterEnd < len(cells) && cells[afterEnd] == game.CellUndetermined {
			cells[afterEnd] = game.CellWhite
			changed = true
		}
	}

	if !changed {
		return nil
	}
	return cells
}
