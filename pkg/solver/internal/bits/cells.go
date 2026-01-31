package bits

import (
	"math/bits"

	"github.com/inahym196/picross-solver/pkg/game"
)

type Bits uint32

func (b Bits) Equals(other Bits) bool { return b == other }

// ValueObject
type Cells struct {
	Len    int
	Blacks Bits
	Whites Bits
}

func FromCells(cells []game.Cell) Cells {
	var blacks, whites Bits
	for i, c := range cells {
		switch c {
		case game.CellBlack:
			blacks |= (1 << i)
		case game.CellWhite:
			whites |= (1 << i)
		}
	}
	return Cells{len(cells), blacks, whites}
}

func NewCells(len int) Cells {
	return Cells{len, 0, 0}
}

func NewCellsWithWhiteMasked(len int) Cells {
	return Cells{len, 0, 1<<len - 1}
}

func NewCellsWithBlackMasked(len int) Cells {
	return Cells{len, 1<<len - 1, 0}
}

func (c Cells) ToCells() []game.Cell {
	cells := make([]game.Cell, c.Len)
	for i := range c.Len {
		mask := Bits(1 << i)
		switch {
		case c.Blacks&mask != 0:
			cells[i] = game.CellBlack
		case c.Whites&mask != 0:
			cells[i] = game.CellWhite
		default:
			cells[i] = game.CellUndetermined
		}
	}
	return cells
}

func (c Cells) Equals(another Cells) bool { return c == another }

//func (c Cells) Mask() Bits { return Bits(1<<c.Len - 1) }

func (c Cells) MarkedBlacks(b Bits) Cells {
	return Cells{c.Len, c.Blacks | b, c.Whites &^ b}
}

func (c Cells) MarkedWhites(b Bits) Cells {
	return Cells{c.Len, c.Blacks &^ b, c.Whites | b}
}

func (c Cells) LeftMostBlackNotWhite() (int, bool) {
	mask := Bits(1<<c.Len - 1)
	target := c.Blacks &^ c.Whites & mask
	if target == 0 {
		return 0, false
	}
	return bits.TrailingZeros32(uint32(target)), true
}
