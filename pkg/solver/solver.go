package solver

import (
	"fmt"
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/rules"
)

type Deduction struct {
	RuleName string
	Hints    []int
	LineRef  game.LineRef
	Before   []game.Cell
	After    []game.Cell
}

func (deduction Deduction) String() string {
	return fmt.Sprintf("%s %s %v %v -> %v", deduction.RuleName, deduction.LineRef, deduction.Hints, deduction.Before, deduction.After)
}

type Rule interface {
	Name() string
	Deduce(game.Line) []game.Cell
}

type Solver struct {
	rules []Rule
}

func NewSolver() *Solver {
	return &Solver{[]Rule{
		rules.ZeroHintRule{},
		rules.MinimumSpacingRule{},
		rules.OverlapFillRule{},
		rules.OverlapExpansionRule{},
		rules.EdgeExpansionRule{},
		rules.MaxHintBlockBoundaryRule{},
		rules.HogeHogeRule{},
		rules.PruneImpossibleSegmentRule{},
		rules.FillRemainingWhiteRule{},
	}}
}

func (s *Solver) ApplyMany(g *game.Game) (int, []Deduction) {
	var ds []Deduction
	for n := 0; ; n++ {
		if dsOnce := s.ApplyOnce(g); len(dsOnce) > 0 {
			ds = append(ds, dsOnce...)
			continue
		}
		return n, ds
	}
}

func (s *Solver) ApplyOnce(g *game.Game) (ds []Deduction) {

	gl := g.Lines()
	for _, l := range gl {
		lds := make([]Deduction, 0)
		for _, rule := range s.rules {
			current := game.Line{Cells: slices.Clone(l.Cells), Hints: l.Hints, Ref: l.Ref}
			if slices.Index(current.Cells, game.CellUndetermined) == -1 {
				break
			}
			before := slices.Clone(current.Cells)
			updated := rule.Deduce(current)

			if updated == nil || slices.Equal(before, updated) {
				continue
			}
			lds = append(lds, Deduction{
				RuleName: rule.Name(),
				Hints:    current.Hints,
				LineRef:  current.Ref,
				Before:   before,
				After:    updated,
			})
		}
		if len(lds) > 0 {
			last := lds[len(lds)-1]
			s.MarkCells(g, last.LineRef, last.After)
			ds = append(ds, lds...)
		}
	}
	return ds
}

func (s *Solver) MarkCells(g *game.Game, ref game.LineRef, cells []game.Cell) {
	switch ref.Kind {
	case game.LineKindRow:
		for i, c := range cells {
			err := g.Mark(ref.Index, i, c)
			if err != nil {
				panic(err)
			}
		}
	case game.LineKindColumn:
		for i, c := range cells {
			err := g.Mark(i, ref.Index, c)
			if err != nil {
				panic(err)
			}
		}
	}
}
