package cells_test

import (
	"testing"

	. "github.com/inahym196/picross-solver/internal/testutil"
	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/cells"
)

func TestGameFromCells(t *testing.T) {
	src := []game.Cell{U, W, B}

	cs := cells.FromGameCells(src)

	if cs.Blacks != 0b100 {
		t.Errorf("expected 0b100, got %#b", cs.Blacks)
	}
	if cs.Whites != 0b10 {
		t.Errorf("expected 0b10, got %#b", cs.Whites)
	}
}
