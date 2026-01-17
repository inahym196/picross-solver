package picrosssolver

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

type Solver struct {
	deducer deducer
}

func NewSolver() Solver {
	return Solver{newDeducer()}
}

func (s Solver) ApplyOnce(game *game.Game) (deds []deduction) {
	board := game.Board()

	for i := range game.RowHints {
		ref := lineRef{lineKindRow, i}
		acc := lineAccessor{&board, ref}
		line := lineView{
			Cells: acc.Cells(),
			Hints: slices.Clone(game.RowHints[i]),
		}

		if lineDeds := s.deducer.DeduceLine(line, ref); len(lineDeds) > 0 {
			last := lineDeds[len(lineDeds)-1]
			acc.Update(last.after)
			deds = append(deds, lineDeds...)
		}
	}
	for i := range game.ColHints {
		ref := lineRef{lineKindColumn, i}
		acc := lineAccessor{&board, ref}
		line := lineView{
			Cells: acc.Cells(),
			Hints: slices.Clone(game.ColHints[i]),
		}

		if lineDeds := s.deducer.DeduceLine(line, ref); len(lineDeds) > 0 {
			last := lineDeds[len(lineDeds)-1]
			acc.Update(last.after)
			deds = append(deds, lineDeds...)
		}
	}
	return deds
}

func (s Solver) ApplyMany(game *game.Game) (int, []deduction) {
	var deds []deduction
	for n := 0; ; n++ {
		if OnceDeds := s.ApplyOnce(game); len(OnceDeds) > 0 {
			deds = append(deds, OnceDeds...)
			continue
		}
		return n, deds
	}
}
