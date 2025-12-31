package picrosssolver

import (
	"slices"
)

type Rule interface {
	Deduce(HintedCells) []Cell
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

// TODO: 実装が雑すぎるので後で綺麗にする。一応テストは通る
func (r OverlapFillRule) Deduce(hc HintedCells) []Cell {
	leftCells := make([]int, len(hc.Cells))
	var last int
	for i, hint := range hc.Hints {
		for range hint {
			leftCells[last] = i + 1
			last++
		}
		if i != len(hc.Hints)-1 {
			last++
		}
	}
	rightCells := make([]int, len(hc.Cells))
	last = 0
	for i, hint := range slices.Backward(hc.Hints) {
		for range hint {
			rightCells[last] = i + 1
			last++
		}
		if i != 0 {
			last++
		}
	}
	slices.Reverse(rightCells)

	cells := make([]Cell, len(hc.Cells))
	copy(cells, hc.Cells)
	for i := range cells {
		if leftCells[i] == rightCells[i] {
			cells[i] = CellBlack
		}
	}
	return cells
}

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpantionRule struct{}

func (r EdgeExpantionRule) countLeading(cells []Cell, cell Cell) int {
	cnt := 0
	for i := range cells {
		if cells[i] != cell {
			return cnt
		}
		cnt++
	}
	return cnt
}

func (r EdgeExpantionRule) lstrip(cells []Cell, cell Cell) []Cell {
	return cells[r.countLeading(cells, cell):]
}

func (r EdgeExpantionRule) applyLeft(cells []Cell, hint int) (changed bool) {
	seg := r.lstrip(cells, CellWhite)
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

func (r EdgeExpantionRule) Deduce(hc HintedCells) []Cell {
	cells := make([]Cell, len(hc.Cells))
	copy(cells, hc.Cells)

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
