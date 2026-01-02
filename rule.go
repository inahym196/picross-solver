package picrosssolver

import (
	"slices"
)

type Rule interface {
	Name() string
	Deduce(HintedCells) []Cell
}

func splitByWhite(cells []Cell) [][]Cell {
	var segs [][]Cell
	var start int
	for i, c := range cells {
		if c == CellWhite {
			if start < i {
				segs = append(segs, cells[start:i])
			}
			start = i + 1
		}
	}
	if start < len(cells) {
		segs = append(segs, cells[start:])
	}
	return segs
}

type ZeroHintRule struct{}

func (e ZeroHintRule) Name() string {
	return "ZeroHintRule"
}

func (r ZeroHintRule) Deduce(hc HintedCells) []Cell {
	if len(hc.Hints) != 1 || hc.Hints[0] != 0 {
		return nil
	}
	for i := range hc.Cells {
		hc.Cells[i] = CellWhite
	}
	return hc.Cells
}

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

func (r MinimumSpacingRule) Name() string {
	return "MinimumSpacingRule"
}

func (r MinimumSpacingRule) Deduce(hc HintedCells) []Cell {

	segs := splitByWhite(hc.Cells)
	if len(segs) != 1 {
		return nil
	}

	cells := segs[0]
	var sum int
	for _, h := range hc.Hints {
		sum += h
	}
	if sum+(len(hc.Hints)-1) != len(cells) {
		return nil
	}

	var last int
	for i, hint := range hc.Hints {
		for range hint {
			cells[last] = CellBlack
			last++
		}
		if i != len(hc.Hints)-1 {
			cells[last] = CellWhite
			last++
		}
	}
	return hc.Cells
}

// ヒントブロックを左詰め／右詰めしたときに必ず重なる部分を黒確定
type OverlapFillRule struct{}

func (r OverlapFillRule) Name() string {
	return "OverlapFillRule"
}

func (r OverlapFillRule) nextPlacablePos(cells []Cell, start int) int {
	for i := start; i < len(cells); i++ {
		if cells[i] != CellWhite {
			return i
		}
	}
	return len(cells)
}

func (r OverlapFillRule) leftAlignedStarts(cells []Cell, hints []int) []int {
	starts := make([]int, len(hints))
	pos := 0

	for i, h := range hints {
		pos = r.nextPlacablePos(cells, pos)
		if pos+h > len(cells) {
			return nil
		}
		for slices.Contains(cells[pos:pos+h], CellWhite) {
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

func (r OverlapFillRule) prevPlacablePos(cells []Cell, start int) int {
	for i := start; i >= 0; i-- {
		if cells[i] != CellWhite {
			return i
		}
	}
	return -1
}

func (r OverlapFillRule) rightAlignedStarts(cells []Cell, hints []int) []int {
	starts := make([]int, len(hints))
	pos := len(cells) - 1

	for i := len(hints) - 1; i >= 0; i-- {
		h := hints[i]

		pos = r.prevPlacablePos(cells, pos)
		start := pos - h + 1
		if start < 0 {
			return nil
		}

		for slices.Contains(cells[start:pos+1], CellWhite) {
			pos := r.prevPlacablePos(cells, start-1)
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

func (r OverlapFillRule) Deduce(hc HintedCells) []Cell {
	cells := hc.Cells

	leftStarts := r.leftAlignedStarts(cells, hc.Hints)
	rightStarts := r.rightAlignedStarts(cells, hc.Hints)

	if leftStarts == nil || rightStarts == nil {
		return nil
	}

	changed := false
	for i, hint := range hc.Hints {
		left := leftStarts[i]
		right := rightStarts[i]

		overlapStart := max(left, right)
		overlapEnd := min(left+hint, right+hint)

		for p := overlapStart; p < overlapEnd; p++ {
			if cells[p] == CellUndetermined {
				cells[p] = CellBlack
				changed = true
			}
		}
	}
	if !changed {
		return nil
	}
	return cells
}

// 端が未確定なら黒をヒント分拡張する
type OverlapExpansionRule struct{}

func (r OverlapExpansionRule) Name() string {
	return "OverlapExpansionRule"
}

func (r OverlapExpansionRule) applyLeft(cells []Cell, hint int) (changed bool) {
	seg := splitByWhite(cells)[0]
	firstBlackIndex := slices.Index(seg, CellBlack)
	if firstBlackIndex == -1 || firstBlackIndex >= hint {
		return false
	}

	for i := firstBlackIndex + 1; i < hint; i++ {
		seg[i] = CellBlack
		changed = true
	}
	return changed
}

func (r OverlapExpansionRule) Deduce(hc HintedCells) []Cell {
	cells := hc.Cells

	firstHint := hc.Hints[0]
	changed1 := r.applyLeft(cells, firstHint)

	slices.Reverse(cells)
	lastHint := hc.Hints[len(hc.Hints)-1]
	changed2 := r.applyLeft(cells, lastHint)

	if !changed1 && !changed2 {
		return nil
	}
	slices.Reverse(cells)
	return cells
}

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) Name() string {
	return "EdgeExpansionRule"
}

func (r EdgeExpansionRule) applyLeft(cells []Cell, hint int) (changed bool) {
	seg := splitByWhite(cells)[0]
	if seg[0] != CellBlack || len(seg) < hint {
		return false
	}
	for i := range hint {
		seg[i] = CellBlack
		changed = true
	}
	if len(seg) > hint {
		seg[hint] = CellWhite
	}
	return changed
}

func (r EdgeExpansionRule) Deduce(hc HintedCells) []Cell {
	cells := hc.Cells

	firstHint := hc.Hints[0]
	changed1 := r.applyLeft(cells, firstHint)

	slices.Reverse(cells)
	lastHint := hc.Hints[len(hc.Hints)-1]
	changed2 := r.applyLeft(cells, lastHint)

	if !changed1 && !changed2 {
		return nil
	}
	slices.Reverse(cells)
	return cells
}

// 既に黒が hint 長に達しているブロックの前後を白確定
type BlockSatisfiedRule struct{}

func (r BlockSatisfiedRule) Name() string {
	return "BlockSatisfiedRule"
}

func (r BlockSatisfiedRule) maxHint(hints []int) int {
	hint := -1
	for _, h := range hints {
		hint = max(hint, h)
	}
	return hint
}

func findSingleBlackBlock(cells []Cell) (start, length int) {
	start = -1
	length = 0

	i := 0
	for i < len(cells) {
		if cells[i] != CellBlack {
			i++
			continue
		}

		if start != -1 {
			return -1, 0
		}

		start = i
		for i < len(cells) && cells[i] == CellBlack {
			length++
			i++
		}
	}

	return start, length
}

func (r BlockSatisfiedRule) Deduce(hc HintedCells) []Cell {

	hint := r.maxHint(hc.Hints)
	cells := hc.Cells

	start, length := findSingleBlackBlock(cells)
	if length != hint {
		return nil
	}

	changed := false

	if start-1 >= 0 && cells[start-1] == CellUndetermined {
		cells[start-1] = CellWhite
		changed = true
	}

	end := start + length
	if end < len(cells) && cells[end] == CellUndetermined {
		cells[end] = CellWhite
		changed = true
	}

	if !changed {
		return nil
	}
	return cells
}

// 最小 hint が収まらない区間を白確定
type PruneImpossibleSegmentRule struct{}

func (r PruneImpossibleSegmentRule) Name() string {
	return "PruneImpossibleSegmentRule"
}

func (r PruneImpossibleSegmentRule) minHint(hints []int) int {
	hint := hints[0]
	for _, h := range hints {
		hint = min(hint, h)
	}
	return hint
}

func (r PruneImpossibleSegmentRule) Deduce(hc HintedCells) []Cell {
	hint := r.minHint(hc.Hints)
	changed := false

	segs := splitByWhite(hc.Cells)
	for i, seg := range segs {
		if len(seg) < hint {
			changed = true
			for j := range seg {
				segs[i][j] = CellWhite
			}
		}
	}
	if !changed {
		return nil
	}
	return hc.Cells
}

// すべての hint を満たした後の残りは白
type FillRemainingWhiteRule struct{}

func (r FillRemainingWhiteRule) Name() string {
	return "FillRemainingWhiteRule"
}

func (r FillRemainingWhiteRule) Deduce(hc HintedCells) []Cell {
	sumHints := 0
	for _, h := range hc.Hints {
		sumHints += h
	}

	blackCount := 0
	for _, c := range hc.Cells {
		if c == CellBlack {
			blackCount++
		}
	}

	if blackCount != sumHints {
		return nil
	}

	deduced := make([]Cell, len(hc.Cells))
	changed := false
	for i, c := range hc.Cells {
		if c == CellUndetermined {
			deduced[i] = CellWhite
			changed = true
		} else {
			deduced[i] = c
		}
	}

	if !changed {
		return nil
	}
	return deduced
}

// 仮に黒／白を置き、矛盾が出たら逆を確定
type HypothesisRule struct{}
