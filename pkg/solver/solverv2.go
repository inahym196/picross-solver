package solver

import (
	math_bits "math/bits"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
	"github.com/inahym196/picross-solver/pkg/solver/internal/history"
	"github.com/inahym196/picross-solver/pkg/solver/internal/queue"
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

type LineRefQueue = queue.Queue[game.LineRef]

type SolverV2 struct {
	rules  []RuleV2
	logger Logger
	q      *LineRefQueue
}

func NewSolverV2(opts ...Option) *SolverV2 {
	s := &SolverV2{[]RuleV2{
		rule.EdgeExpansionRule{},
		rule.ExactMatchRule{},
	}, nopLogger{}, &LineRefQueue{}}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *SolverV2) Solve(g *game.Game, h *history.History) {
	diffs := s.initProject(g)

	for _, diff := range diffs {
		before := g.Line(diff.Ref)
		s.markCells(g, diff.Ref, diff.Cells)
		s.logger.Logf("projected: %v %v -> %v", before, diff.Cells, g.Line(diff.Ref).Cells)
		for _, ref := range diff.RefsInDiff() {
			s.q.Push(ref)
		}
	}

	for {
		ref, ok := s.q.Pop()
		if !ok {
			break
		}
		line := g.Line(ref)
		if line.IsFilled() {
			continue
		}
		s.logger.Logf("start line:%v", line)
		d, err := domain.NewLineDomain(line.Len(), line.Hints)
		if err != nil {
			panic(err)
		}
		current := bits.FromCells(line.Cells)
		last, narrowed := s.narrowLine(current, d, h)
		if !narrowed {
			continue
		}
		projected, err := last.Project()
		if err != nil {
			panic(err)
		}
		diff := NewLineDiff(ref, current, projected)
		if !diff.IsEmpty() {
			s.markCells(g, ref, diff.Cells)
			s.logger.Logf("projected: %v %v -> %v", current, diff.Cells, g.Line(ref).Cells)
			for _, next := range diff.RefsInDiff() {
				s.q.Push(next)
			}
		}
	}
}

type LineDiff struct {
	Ref   game.LineRef
	Cells bits.Cells
}

func NewLineDiff(ref game.LineRef, before bits.Cells, after bits.Cells) LineDiff {
	return LineDiff{ref, after.ExtraCellsFrom(before)}
}

func (diff LineDiff) IsEmpty() bool { return diff.Cells.IsEmpty() }

func (diff LineDiff) RefsInDiff() []game.LineRef {
	cells := diff.Cells
	if cells.IsEmpty() {
		return nil
	}

	mask := uint32(cells.Blacks | cells.Whites)
	refs := make([]game.LineRef, 0, math_bits.OnesCount32(mask))
	for mask != 0 {
		i := math_bits.TrailingZeros32(mask)
		refs = append(refs, game.LineRef{Kind: diff.Ref.Kind, Index: i})
		mask &= mask - 1
	}
	return refs
}

func (s *SolverV2) initProject(g *game.Game) []LineDiff {
	diffs := make([]LineDiff, 0, g.Width()+g.Height())

	for row := range g.AllRows() {
		diff := NewLineDiff(row.Ref, bits.FromCells(row.Cells), s.projectLine(row))
		if !diff.IsEmpty() {
			diffs = append(diffs, diff)
		}
	}

	for col := range g.AllColumns() {
		diff := NewLineDiff(col.Ref, bits.FromCells(col.Cells), s.projectLine(col))
		if !diff.IsEmpty() {
			diffs = append(diffs, diff)
		}
	}
	return diffs
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
			changed, err := g.MarkByRef(ref, i, cell)
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
