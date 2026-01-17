package picrosssolver

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

func splitByWhite(cells []game.Cell) [][]game.Cell {
	var segs [][]game.Cell
	var start int
	for i, c := range cells {
		if c == game.CellWhite {
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

func (r ZeroHintRule) Deduce(line lineView) []game.Cell {
	cells := slices.Clone(line.Cells)
	if len(line.Hints) != 1 || line.Hints[0] != 0 {
		return nil
	}
	for i := range cells {
		cells[i] = game.CellWhite
	}
	return cells
}

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

func (r MinimumSpacingRule) Name() string {
	return "MinimumSpacingRule"
}

func (r MinimumSpacingRule) Deduce(line lineView) []game.Cell {
	cells := slices.Clone(line.Cells)

	segs := splitByWhite(cells)
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

func (r OverlapFillRule) Deduce(line lineView) []game.Cell {
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

// 端が未確定なら黒をヒント分拡張する
type OverlapExpansionRule struct{}

func (r OverlapExpansionRule) Name() string {
	return "OverlapExpansionRule"
}

func (r OverlapExpansionRule) applyLeft(cells []game.Cell, hint int) (changed bool) {
	seg := splitByWhite(cells)[0]
	firstBlackIndex := slices.Index(seg, game.CellBlack)
	if firstBlackIndex == -1 || firstBlackIndex >= hint {
		return false
	}

	for i := firstBlackIndex + 1; i < hint; i++ {
		// TODO: バグの可能性あり
		seg[i] = game.CellBlack
		changed = true
	}
	return changed
}

func (r OverlapExpansionRule) Deduce(line lineView) []game.Cell {
	cells := slices.Clone(line.Cells)

	firstHint := line.Hints[0]
	changed1 := r.applyLeft(cells, firstHint)

	slices.Reverse(cells)
	lastHint := line.Hints[len(line.Hints)-1]
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

func (r EdgeExpansionRule) applyLeft(cells []game.Cell, hint int) (changed bool) {
	seg := splitByWhite(cells)[0]
	if seg[0] != game.CellBlack || len(seg) < hint {
		return false
	}
	for i := range hint {
		seg[i] = game.CellBlack
		changed = true
	}
	if len(seg) > hint {
		seg[hint] = game.CellWhite
	}
	return changed
}

func (r EdgeExpansionRule) Deduce(line lineView) []game.Cell {
	cells := slices.Clone(line.Cells)

	firstHint := line.Hints[0]
	changed1 := r.applyLeft(cells, firstHint)

	slices.Reverse(cells)
	lastHint := line.Hints[len(line.Hints)-1]
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

type Block struct {
	start  int
	length int
}

func nextBlock(cells []game.Cell, start int) *Block {
	for cells[start] != game.CellBlack {
		start++
		if start >= len(cells) {
			return nil
		}
	}
	length := 0
	end := start
	for end < len(cells) && cells[end] == game.CellBlack {
		end++
		length++
	}
	return &Block{start, length}
}

func findBlocksN(cells []game.Cell, n int) []Block {
	var blocks []Block
	for i := 0; i < len(cells); i++ {
		block := nextBlock(cells, i)
		if block == nil {
			return blocks
		}
		if block.length != n {
			continue
		}
		blocks = append(blocks, *block)
		i = block.start + block.length + 1
	}
	return blocks
}

func (r BlockSatisfiedRule) Deduce(line lineView) []game.Cell {
	cells := slices.Clone(line.Cells)
	hint := r.maxHint(line.Hints)
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

func (r PruneImpossibleSegmentRule) Deduce(line lineView) []game.Cell {
	cells := slices.Clone(line.Cells)

	hint := r.minHint(line.Hints)
	changed := false

	segs := splitByWhite(cells)
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

// すべての hint を満たした後の残りは白
type FillRemainingWhiteRule struct{}

func (r FillRemainingWhiteRule) Name() string {
	return "FillRemainingWhiteRule"
}

func (r FillRemainingWhiteRule) Deduce(line lineView) []game.Cell {
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

type hogehogeRule struct{}

func (r hogehogeRule) Name() string {
	return "hogehogeRule"
}

// 左端から広義単調増加列を取得
func weaklyIncreasingFromLeft(a []int) []int {
	if len(a) == 0 {
		return nil
	}

	res := []int{a[0]}
	for i := 1; i < len(a); i++ {
		if a[i] < a[i-1] {
			break
		}
		res = append(res, a[i])
	}
	return res
}

// 右端から広義単調増加列を取得
func weaklyIncreasingFromRight(a []int) []int {
	if len(a) == 0 {
		return nil
	}

	tmp := []int{a[len(a)-1]}
	for i := len(a) - 2; i >= 0; i-- {
		if a[i] < a[i+1] {
			break
		}
		tmp = append(tmp, a[i])
	}

	slices.Reverse(tmp)
	return tmp
}

func (r hogehogeRule) nextPlacablePos(cells []game.Cell, start int) int {
	for i := start; i < len(cells); i++ {
		if cells[i] != game.CellWhite {
			return i
		}
	}
	return len(cells)
}

func (r hogehogeRule) leftAlignedEnds(cells []game.Cell, hints []int) []int {
	ends := make([]int, len(hints))
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
		ends[i] = pos + h - 1
		pos += h + 1
	}
	return ends
}

func (r hogehogeRule) prevPlacablePos(cells []game.Cell, start int) int {
	for i := start; i >= 0; i-- {
		if cells[i] != game.CellWhite {
			return i
		}
	}
	return -1
}

func (r hogehogeRule) rightAlignedStarts(cells []game.Cell, hints []int) []int {
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

func (r hogehogeRule) getBlockAt(cells []game.Cell, index int) (Block, bool) {
	if index < 0 || index >= len(cells) {
		return Block{}, false
	}
	if cells[index] != game.CellBlack {
		return Block{}, false
	}
	start := index
	for start > 0 && cells[start-1] == game.CellBlack {
		start--
	}

	end := index
	for end+1 < len(cells) && cells[end+1] == game.CellBlack {
		end++
	}
	return Block{start, end - start + 1}, true
}

func (r hogehogeRule) tryDrawWhite(cells []game.Cell, index int) bool {
	if 0 <= index && index < len(cells) {
		cells[index] = game.CellWhite
		return true
	}
	return false
}

func (r hogehogeRule) Deduce(line lineView) []game.Cell {
	cells := slices.Clone(line.Cells)
	changed := false

	if hints := weaklyIncreasingFromLeft(line.Hints); len(hints) > 0 {
		for i, end := range r.leftAlignedEnds(cells, hints) {
			block, found := r.getBlockAt(cells, end)
			if !found {
				block, found = r.getBlockAt(cells, end+1)
				if !found {
					continue
				}
			}
			if hints[i] == block.length {
				if r.tryDrawWhite(cells, block.start-1) {
					changed = true
				}
				if r.tryDrawWhite(cells, block.start+block.length) {
					changed = true
				}
			}
		}
	}

	if hints := weaklyIncreasingFromRight(line.Hints); len(hints) > 0 {
		for i, start := range r.rightAlignedStarts(line.Cells, hints) {
			block, found := r.getBlockAt(cells, start)
			if !found {
				block, found = r.getBlockAt(cells, start-1)
				if !found {
					continue
				}
			}
			if hints[i] == block.length {
				if r.tryDrawWhite(cells, block.start-1) {
					changed = true
				}
				if r.tryDrawWhite(cells, block.start+block.length) {
					changed = true
				}
			}
		}
	}
	if !changed {
		return nil
	}
	return cells
}

// 仮に黒／白を置き、矛盾が出たら逆を確定
type HypothesisRule struct{}
