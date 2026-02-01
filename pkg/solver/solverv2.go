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

func (s *SolverV2) ApplyMany(g *game.Game) (n int, h *history.History) {
	return 1, s.Apply(g)
}

func (s *SolverV2) Apply(g *game.Game) (h *history.History) {

	for row := range g.AllRows() {
		d, err := domain.NewLineDomain(g.Width(), row.Hints)
		if err != nil {
			panic(err)
		}
		current := bits.FromCells(row.Cells)
		if d.IsDeterministic() {
			updated, err := d.Project()
			if err != nil {
				panic(err)
			}
			s.applyProjection(g, row.Ref, current, updated)
			continue
		}
		lh := s.NarrowLine(bits.FromCells(row.Cells), d)
		if lh.IsEmpty() {
			continue
		}
		updated, err := lh.Last().Domain.Project()
		if err != nil {
			panic(err)
		}
		s.applyProjection(g, row.Ref, current, updated)
		h.Merge(lh)
	}

	for col := range g.AllColumns() {
		d, err := domain.NewLineDomain(g.Width(), col.Hints)
		if err != nil {
			panic(err)
		}

		current := bits.FromCells(col.Cells)
		if d.IsDeterministic() {
			pj, err := d.Project()
			if err != nil {
				panic(err)
			}
			s.applyProjection(g, col.Ref, current, pj)
			continue
		}

		lh := s.NarrowLine(bits.FromCells(col.Cells), d)
		if lh.IsEmpty() {
			continue
		}

		updated, err := lh.Last().Domain.Project()
		if err != nil {
			panic(err)
		}
		s.applyProjection(g, col.Ref, current, updated)
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
