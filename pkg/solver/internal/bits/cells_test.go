package bits_test

import (
	"testing"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
)

const (
	U = game.CellUndetermined
	W = game.CellWhite
	B = game.CellBlack
)

func TestGameFromCells(t *testing.T) {
	src := []game.Cell{U, W, B}

	cells := bits.FromCells(src)

	if cells.Blacks != 0b100 {
		t.Errorf("expected 0b100, got %#b", cells.Blacks)
	}
	if cells.Whites != 0b10 {
		t.Errorf("expected 0b10, got %#b", cells.Whites)
	}
}
