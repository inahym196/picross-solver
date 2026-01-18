package line

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

type Line struct {
	Cells []game.Cell
	Hints []int
}

func (l Line) IsFilled() bool {
	return slices.Index(l.Cells, game.CellUndetermined) == -1
}
