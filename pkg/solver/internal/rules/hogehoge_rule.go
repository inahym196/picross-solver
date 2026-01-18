package rules

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

type HogeHogeRule struct{}

func (r HogeHogeRule) Name() string {
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

func (r HogeHogeRule) nextPlacablePos(cells []game.Cell, start int) int {
	for i := start; i < len(cells); i++ {
		if cells[i] != game.CellWhite {
			return i
		}
	}
	return len(cells)
}

func (r HogeHogeRule) leftAlignedEnds(cells []game.Cell, hints []int) []int {
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

func (r HogeHogeRule) prevPlacablePos(cells []game.Cell, start int) int {
	for i := start; i >= 0; i-- {
		if cells[i] != game.CellWhite {
			return i
		}
	}
	return -1
}

func (r HogeHogeRule) rightAlignedStarts(cells []game.Cell, hints []int) []int {
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

func (r HogeHogeRule) getBlockAt(cells []game.Cell, index int) (Block, bool) {
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

func (r HogeHogeRule) tryDrawWhite(cells []game.Cell, index int) bool {
	if 0 <= index && index < len(cells) {
		cells[index] = game.CellWhite
		return true
	}
	return false
}

func (r HogeHogeRule) Deduce(line line.Line) []game.Cell {
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
