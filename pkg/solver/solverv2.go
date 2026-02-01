package solver

import (
	math_bits "math/bits"

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

func (s *SolverV2) ApplyMany(g *game.Game, h *history.History) (n int) {
	s.Apply(g, h)
	return 1
}

func (s *SolverV2) Apply(g *game.Game, h *history.History) {
	for row := range g.AllRows() {
		s.applyLine(g, row, h)
	}
	for col := range g.AllColumns() {
		s.applyLine(g, col, h)
	}
}

func (s *SolverV2) applyLine(g *game.Game, l game.Line, h *history.History) {
	d, err := domain.NewLineDomain(g.Width(), l.Hints)
	if err != nil {
		panic(err)
	}
	current := bits.FromCells(l.Cells)
	if d.IsDeterministic() {
		updated, err := d.Project()
		if err != nil {
			panic(err)
		}
		s.applyProjection(g, l.Ref, current, updated)
		return
	}
	lastD, narrowed := s.narrowLine(bits.FromCells(l.Cells), d, h)
	if narrowed == false {
		return
	}
	updated, err := lastD.Project()
	if err != nil {
		panic(err)
	}
	s.applyProjection(g, l.Ref, current, updated)
}

func (s *SolverV2) narrowLine(cells bits.Cells, d domain.LineDomain, h *history.History) (last domain.LineDomain, narrowed bool) {
	for _, rule := range s.rules {
		newD, changed := rule.Narrow(cells, d)
		if !changed || newD.Equals(d) {
			continue
		}
		d = newD
		h.Append(history.Step{RuleName: rule.Name(), Domain: newD})
		narrowed = true
	}
	return d, narrowed
}

func (s *SolverV2) applyProjection(g *game.Game, ref game.LineRef, current bits.Cells, projected bits.Cells) {
	updated, conflict := current.Merged(projected)
	if conflict {
		panic("conflict")
	}
	if updated.Equals(current) {
		return
	}
	s.MarkCells(g, ref, updated)
}

func (s *SolverV2) MarkCells(g *game.Game, ref game.LineRef, cells bits.Cells) {

	processBits := func(bits uint32, cell game.Cell) {
		for bits != 0 {
			i := math_bits.TrailingZeros32(bits)
			var row, col int
			if ref.Kind == game.LineKindRow {
				row, col = ref.Index, i
			} else {
				row, col = i, ref.Index
			}
			g.Mark(row, col, cell)
			bits &= bits - 1
		}
	}
	processBits(uint32(cells.Blacks), game.CellBlack)
	processBits(uint32(cells.Whites), game.CellWhite)
}
