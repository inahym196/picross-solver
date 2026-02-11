package rule

import (
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
)

type ExactMatchRule struct{}

func (r ExactMatchRule) Name() string { return "ExactMatchRule" }

func (r ExactMatchRule) Narrow(cells bits.Cells, d domain.LineDomain) (domain.LineDomain, bool) {
	hints := 0
	for _, run := range d.AllRuns() {
		hints += run.Len
	}
	if hints != cells.Blacks.OnesCount() {
		return d, false
	}
	cursor := 0
	changed := false
	for i, run := range d.AllRuns() {
		start := cursor + (cells.Blacks>>cursor).LeftZeros()
		run = run.Fixed(start)

		var narrowed bool
		d, narrowed = d.Narrowed(i, run)
		changed = changed || narrowed

		cursor = start + run.Len
	}
	return d, changed
}
