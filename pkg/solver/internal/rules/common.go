package rules

import (
	"github.com/inahym196/picross-solver/pkg/game"
)

func SplitByWhite(cells []game.Cell) [][]game.Cell {
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

// 仮に黒／白を置き、矛盾が出たら逆を確定
type HypothesisRule struct{}
