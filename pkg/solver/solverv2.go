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

type Logger interface {
	Logf(format string, args ...any)
	Verbose() bool
}

type nopLogger struct{}

func (nopLogger) Logf(string, ...any) {}
func (nopLogger) Verbose() bool       { return false }

type Option func(*SolverV2)

func WithLogger(l Logger) Option {
	return func(s *SolverV2) {
		s.logger = l
	}
}

type SolverV2 struct {
	rules  []RuleV2
	logger Logger
}

func NewSolverV2(opts ...Option) *SolverV2 {
	s := &SolverV2{[]RuleV2{rule.EdgeExpansionRule{}}, nopLogger{}}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *SolverV2) ApplyMany(g *game.Game, h *history.History) (n int) {
	for row := range g.AllRows() {
		diffs := s.projectLine(row).ExtraCellsFrom(bits.FromCells(row.Cells))
		if diffs.IsEmpty() {
			continue
		}
		s.logger.Logf("init projected: %v %v", row, diffs)
		s.markCells(g, row.Ref, diffs)
	}

	for col := range g.AllColumns() {
		diffs := s.projectLine(col).ExtraCellsFrom(bits.FromCells(col.Cells))
		if diffs.IsEmpty() {
			continue
		}
		s.logger.Logf("init projected: %v %v", col, diffs)
		s.markCells(g, col.Ref, diffs)
	}

	for row := range g.AllRows() {
		s.applyLine(g, row, h)
	}
	for col := range g.AllColumns() {
		s.applyLine(g, col, h)
	}
	return 1
}

func (s *SolverV2) projectLine(l game.Line) bits.Cells {
	d, err := domain.NewLineDomain(l.Len(), l.Hints)
	if err != nil {
		panic(err)
	}
	projected, err := d.Project()
	if err != nil {
		panic(err)
	}
	return projected
}

func (s *SolverV2) applyLine(g *game.Game, l game.Line, h *history.History) {
	s.logger.Logf("start line:%v", l)
	d, err := domain.NewLineDomain(g.Width(), l.Hints)
	if err != nil {
		panic(err)
	}
	current := bits.FromCells(l.Cells)
	lastD, narrowed := s.narrowLine(current, d, h)
	if !narrowed {
		return
	}
	updated, err := lastD.Project()
	if err != nil {
		panic(err)
	}
	changed := s.markCells(g, l.Ref, updated)
	if changed {
		s.logger.Logf("cells updated: -> %v", updated)
	}
}

func (s *SolverV2) narrowLine(cells bits.Cells, d domain.LineDomain, h *history.History) (last domain.LineDomain, narrowed bool) {
	for _, rule := range s.rules {
		newD, changed := rule.Narrow(cells, d)
		if !changed || newD.Equals(d) {
			if s.logger.Verbose() {
				s.logger.Logf("%s: nochange", rule.Name())
			}
			continue
		}

		s.logger.Logf("%s: narrowed: %v -> %v", rule.Name(), d, newD)
		d = newD
		h.Append(history.Step{RuleName: rule.Name(), Domain: newD})
		narrowed = true
	}
	return d, narrowed
}

func (s *SolverV2) markCells(g *game.Game, ref game.LineRef, cells bits.Cells) bool {
	changedAny := false
	processBits := func(bits uint32, cell game.Cell) {
		for bits != 0 {
			i := math_bits.TrailingZeros32(bits)
			var row, col int
			if ref.Kind == game.LineKindRow {
				row, col = ref.Index, i
			} else {
				row, col = i, ref.Index
			}
			changed, err := g.Mark(row, col, cell)
			if err != nil {
				panic(err)
			}
			if changed {
				changedAny = true
			}
			bits &= bits - 1
		}
	}
	processBits(uint32(cells.Blacks), game.CellBlack)
	processBits(uint32(cells.Whites), game.CellWhite)
	return changedAny
}
