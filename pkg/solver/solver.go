package solver

import (
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/accessor"
	"github.com/inahym196/picross-solver/pkg/solver/internal/deducer"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

type Solver struct {
	deducer deducer.Deducer
}

func NewSolver() Solver {
	return Solver{deducer.NewDeducer()}
}

func (s Solver) ApplyOnce(g game.Game) (deds []deducer.Deduction) {

	for i := range g.RowHints {
		acc := accessor.NewLineAccessor(g, game.LineKindRow, i)
		line := line.Line{
			Cells: acc.Cells(),
			Hints: slices.Clone(g.RowHints[i]),
		}

		if lineDeds := s.deducer.DeduceLine(line, acc.Ref()); len(lineDeds) > 0 {
			last := lineDeds[len(lineDeds)-1]
			acc.Update(last.After)
			deds = append(deds, lineDeds...)
		}
	}
	for i := range g.ColHints {
		acc := accessor.NewLineAccessor(g, game.LineKindColumn, i)
		line := line.Line{
			Cells: acc.Cells(),
			Hints: slices.Clone(g.ColHints[i]),
		}
		if lineDeds := s.deducer.DeduceLine(line, acc.Ref()); len(lineDeds) > 0 {
			last := lineDeds[len(lineDeds)-1]
			acc.Update(last.After)
			deds = append(deds, lineDeds...)
		}
	}
	return deds
}

func (s Solver) ApplyMany(game game.Game) (int, []deducer.Deduction) {
	var deds []deducer.Deduction
	for n := 0; ; n++ {
		if OnceDeds := s.ApplyOnce(game); len(OnceDeds) > 0 {
			deds = append(deds, OnceDeds...)
			continue
		}
		return n, deds
	}
}
