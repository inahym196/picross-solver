package cells

import (
	"github.com/inahym196/picross-solver/pkg/game"
)

type BitSet uint32

type BitCells struct {
	Blacks BitSet
	Whites BitSet
	Len    int
}

func FromGameCells(cells []game.Cell) BitCells {
	var blacks, whites BitSet
	for i, c := range cells {
		switch c {
		case game.CellBlack:
			blacks |= (1 << i)
		case game.CellWhite:
			whites |= (1 << i)
		}
	}
	return BitCells{blacks, whites, len(cells)}
}
