package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

// 最小 hint が収まらない区間を白確定
type PruneImpossibleSegmentRule struct{}

func (r PruneImpossibleSegmentRule) Name() string {
	return "PruneImpossibleSegmentRule"
}

func (r PruneImpossibleSegmentRule) Deduce(line line.Line) []game.Cell {
	cells := slices.Clone(line.Cells)

	hint := slices.Min(line.Hints)
	changed := false

	segs := SplitByWhite(cells)
	for i, seg := range segs {
		if len(seg) < hint {
			changed = true
			for j := range seg {
				segs[i][j] = game.CellWhite
			}
		}
	}
	if !changed {
		return nil
	}
	return cells
}
