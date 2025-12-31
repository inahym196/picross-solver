package picrosssolver

import (
	"slices"
)

type Rule interface {
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

type ExtractMatchRule struct{}

func (r ExtractMatchRule) Deduce(hc HintedCells) []Cell {
	if len(hc.Hints) == 1 && hc.Hints[0] == len(hc.Cells) {
		return filledCells(len(hc.Cells), CellBlack)
	}
	return nil
}

type ZeroHintRule struct{}

func (r ZeroHintRule) Deduce(hc HintedCells) []Cell {
	if len(hc.Hints) == 1 && hc.Hints[0] == 0 {
		return filledCells(len(hc.Cells), CellWhite)
	}
	return nil
}

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

func (r MinimumSpacingRule) Deduce(hc HintedCells) []Cell {

	// 判定部分
	var sum int
	for _, h := range hc.Hints {
		sum += h
	}
	if sum+(len(hc.Hints)-1) != len(hc.Cells) {
		return nil
	}

	// 生成部分
	deduced := make([]Cell, len(hc.Cells))
	var last int
	for i, hint := range hc.Hints {
		for range hint {
			deduced[last] = CellBlack
			last++
		}
		if i != len(hc.Hints)-1 {
			deduced[last] = CellWhite
			last++
		}
	}
	return deduced
}

// ヒントブロックを左詰め／右詰めしたときに必ず重なる部分を黒確定
type OverlapFillRule struct{}

func (r OverlapFillRule) leftAlignedStarts(hints []int) []int {
	starts := make([]int, len(hints))
	pos := 0
	for i, h := range hints {
		starts[i] = pos
		pos += h + 1
	}
	return starts
}

func (r OverlapFillRule) rightAlignedStarts(hints []int, length int) []int {
	starts := make([]int, len(hints))
	pos := length
	for i := len(hints) - 1; i >= 0; i-- {
		pos -= hints[i]
		starts[i] = pos
		pos--
	}
	return starts
}

func (r OverlapFillRule) Deduce(hc HintedCells) []Cell {
	n := len(hc.Cells)
	cells := hc.Cells

	leftStarts := r.leftAlignedStarts(hc.Hints)
	rightStarts := r.rightAlignedStarts(hc.Hints, n)

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

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) applyLeft(cells []Cell, hint int) (changed bool) {
	seg := splitByWhite(cells)[0]
	firstBlackIndex := slices.Index(seg, CellBlack)
	if firstBlackIndex == -1 || firstBlackIndex >= hint {
		return false
	}

	expanding := false
	for i := firstBlackIndex; i < hint && i < len(seg); i++ {
		if seg[i] == CellBlack {
			expanding = true
		} else if expanding {
			seg[i] = CellBlack
			changed = true
		}
	}
	return changed
}

func (r EdgeExpansionRule) Deduce(hc HintedCells) []Cell {
	cells := hc.Cells

	firstHint := hc.Hints[0]
	r.applyLeft(cells, firstHint)
	slices.Reverse(cells)

	lastHint := hc.Hints[len(hc.Hints)-1]
	r.applyLeft(cells, lastHint)
	slices.Reverse(cells)

	return cells
}

// 既に黒が hint 長に達しているブロックの前後を白確定
type BlockSatisfiedRule struct{}

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

// 白確定セルで line を分割し、それぞれにヒントを再配分
type SegmentSplitRule struct{}

// ヒントが収まらない区間を白確定
type PruneImpossibleSegmentRule struct{}

// 長さ < 最小 hint の区間はすべて白
type TooSmallSegmentRule struct{}

// すべての hint を満たした後の残りは白
type FillRemainingWhiteRule struct{}

// 仮に黒／白を置き、矛盾が出たら逆を確定
type HypothesisRule struct{}
