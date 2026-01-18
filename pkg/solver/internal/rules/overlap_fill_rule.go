package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

// ヒントブロックを左詰め／右詰めしたときに必ず重なる部分を黒確定
type OverlapFillRule struct{}

func (r OverlapFillRule) Name() string {
	return "OverlapFillRule"
}

func (r OverlapFillRule) nextPlacablePos(cells []game.Cell, start int) int {
	for i := start; i < len(cells); i++ {
		if cells[i] != game.CellWhite {
			return i
		}
	}
	return len(cells)
}

func (r OverlapFillRule) leftAlignedStarts(cells []game.Cell, hints []int) []int {
	starts := make([]int, len(hints))
	pos := 0

	for i, h := range hints {
		pos = r.nextPlacablePos(cells, pos)
		if pos+h > len(cells) {
			return nil
		}
		for slices.Contains(cells[pos:pos+h], game.CellWhite) {
			pos = r.nextPlacablePos(cells, pos+1)
			if pos+h > len(cells) {
				return nil
			}
		}
		starts[i] = pos
		pos += h + 1
	}
	return starts
}

func (r OverlapFillRule) prevPlacablePos(cells []game.Cell, start int) int {
	for i := start; i >= 0; i-- {
		if cells[i] != game.CellWhite {
			return i
		}
	}
	return -1
}

func (r OverlapFillRule) rightAlignedStarts(cells []game.Cell, hints []int) []int {
	starts := make([]int, len(hints))
	pos := len(cells) - 1
	for i := len(hints) - 1; i >= 0; i-- {
		h := hints[i]

		pos = r.prevPlacablePos(cells, pos)
		start := pos - h + 1
		if start < 0 {
			return nil
		}

		for slices.Contains(cells[start:pos+1], game.CellWhite) {
			pos = r.prevPlacablePos(cells, pos-1)
			start = pos - h + 1
			if start < 0 {
				return nil
			}
		}
		starts[i] = start
		pos = start - 2
	}
	return starts
}

func (r OverlapFillRule) Deduce(line line.Line) []game.Cell {
	cells := slices.Clone(line.Cells)

	leftStarts := r.leftAlignedStarts(cells, line.Hints)
	rightStarts := r.rightAlignedStarts(cells, line.Hints)

	if leftStarts == nil || rightStarts == nil {
		return nil
	}

	changed := false
	for i, hint := range line.Hints {
		left := leftStarts[i]
		right := rightStarts[i]

		overlapStart := max(left, right)
		overlapEnd := min(left+hint, right+hint)

		for p := overlapStart; p < overlapEnd; p++ {
			if cells[p] == game.CellUndetermined {
				cells[p] = game.CellBlack
				changed = true
			}
		}
	}
	if !changed {
		return nil
	}
	return cells
}
