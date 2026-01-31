package solver

import (
	"fmt"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
	"github.com/inahym196/picross-solver/pkg/solver/internal/history"
	"github.com/inahym196/picross-solver/pkg/solver/internal/rule"
)

type RuleV2 interface {
	Name() string
	Narrow(cells bits.Cells, domain domain.LineDomain) (domain.LineDomain, bool)
}

type SolverV2 struct {
	rules []RuleV2
}

func NewSolverV2() *SolverV2 {
	return &SolverV2{[]RuleV2{
		rule.EdgeExpansionRule{},
	}}
}

func (s *SolverV2) ApplyMany(g *game.Game) (n int, h *history.History) {
	return 1, s.Apply(g)
}

func (s *SolverV2) Apply(g *game.Game) (h *history.History) {
	for _, gl := range g.Lines() {
		domain := domain.NewLineDomain(g.Width(), gl.Hints)
		lh := s.NarrowLine(bits.FromCells(gl.Cells), domain)
		if lh.IsEmpty() {
			continue
		}
		cells := lh.Last().Domain.Project()
		fmt.Printf("cells: %v\n", cells)
		s.MarkCells(g, gl.Ref, cells)
		h.Merge(lh)
	}
	return h
}

func (s *SolverV2) NarrowLine(cells bits.Cells, d domain.LineDomain) (h history.History) {
	for _, rule := range s.rules {
		newD, changed := rule.Narrow(cells, d)
		if !changed || newD.Equals(d) {
			continue
		}
		d = newD
		h.Append(history.Step{RuleName: rule.Name(), Domain: newD})
	}
	return h
}

func (s *SolverV2) MarkCells(g *game.Game, ref game.LineRef, cells bits.Cells) {
	//panic("not implemented yet")
}
