package picrosssolver

import (
	"slices"
)

type Solver struct {
	deducer deducer
}

func NewSolver() Solver {
	return Solver{newDeducer()}
}

func (s Solver) ApplyOnce(game *Game) (deds []deduction) {

	for i := range game.rowHints {
		ref := lineRef{lineKindRow, i}
		acc := lineAccessor{&game.board, ref}
		line := lineView{
			Cells: acc.Cells(),
			Hints: slices.Clone(game.rowHints[i]),
		}

		if lineDeds := s.deducer.DeduceLine(line, ref); len(lineDeds) > 0 {
			last := lineDeds[len(lineDeds)-1]
			acc.Update(last.after)
			deds = append(deds, lineDeds...)
		}
	}
	for i := range game.colHints {
		ref := lineRef{lineKindColumn, i}
		acc := lineAccessor{&game.board, ref}
		line := lineView{
			Cells: acc.Cells(),
			Hints: slices.Clone(game.colHints[i]),
		}

		if lineDeds := s.deducer.DeduceLine(line, ref); len(lineDeds) > 0 {
			last := lineDeds[len(lineDeds)-1]
			acc.Update(last.after)
			deds = append(deds, lineDeds...)
		}
	}
	return deds
}

func (s Solver) ApplyMany(game *Game) (int, []deduction) {
	var deds []deduction
	for n := 0; ; n++ {
		if OnceDeds := s.ApplyOnce(game); len(OnceDeds) > 0 {
			deds = append(deds, OnceDeds...)
			continue
		}
		return n, deds
	}
}
