package rule

import (
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
)

type ExactMatchRule struct{}

type ExactMatchBlackRule struct{}

type ExactMatchWhiteRule struct{}

func (r ExactMatchBlackRule) Name() string { return "ExactMatchBlackRule" }

func (r ExactMatchWhiteRule) Name() string { return "ExactMatchWhiteRule" }

func (r ExactMatchBlackRule) Narrow(cells bits.Cells, d domain.LineDomain) (domain.LineDomain, bool) {
	hints := 0
	for _, run := range d.AllRuns() {
		hints += run.Len
	}
	mask := bits.Bits(1<<cells.Len - 1)
	nonWhites := mask &^ cells.Whites
	if hints != nonWhites.OnesCount() {
		return d, false
	}
	return d.FixedByMask(nonWhites)
}

func (r ExactMatchWhiteRule) Narrow(cells bits.Cells, d domain.LineDomain) (domain.LineDomain, bool) {
	hints := 0
	for _, run := range d.AllRuns() {
		hints += run.Len
	}
	if hints != cells.Blacks.OnesCount() {
		return d, false
	}
	return d.FixedByMask(cells.Blacks)
}
