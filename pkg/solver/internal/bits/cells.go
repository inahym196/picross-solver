package bits

import (
	"fmt"
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

func FromMasks(len int, blacks, whites Bits) (Cells, error) {
	var err error
	c := NewCells(len)

	c, err = c.MarkedBlacks(blacks)
	if err != nil {
		return Cells{}, err
	}

	c, err = c.MarkedWhites(whites)
	if err != nil {
		return Cells{}, err
	}

	return c, nil
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

func (c Cells) String() string {
	return fmt.Sprintf("{Len:%d Blacks:b%b Whites:b%b}", c.Len, c.Blacks, c.Whites)
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

func (c Cells) MarkedBlacks(b Bits) (Cells, error) {
	mask := Bits(1<<c.Len - 1)
	if b&^mask != 0 {
		return Cells{}, fmt.Errorf("out of range")
	}

	if b&c.Whites != 0 {
		return Cells{}, fmt.Errorf("conflict")
	}

	return Cells{c.Len, c.Blacks | b, c.Whites}, nil
}

func (c Cells) MarkedWhites(b Bits) (Cells, error) {
	mask := Bits(1<<c.Len - 1)
	if b&^mask != 0 {
		return Cells{}, fmt.Errorf("out of range")
	}

	if b&c.Blacks != 0 {
		return Cells{}, fmt.Errorf("conflict")
	}

	return Cells{c.Len, c.Blacks, c.Whites | b}, nil
}

func (c Cells) MostLeftBlack() int {
	return bits.TrailingZeros32(uint32(c.Blacks))
}

func (c Cells) Merged(other Cells) (Cells, bool) {
	if c.Len != other.Len {
		panic("cells length mismatch")
	}
	blacks := c.Blacks | other.Blacks
	whites := c.Whites | other.Whites
	if blacks&whites != 0 {
		return Cells{}, true
	}
	return Cells{c.Len, blacks, whites}, false
}
