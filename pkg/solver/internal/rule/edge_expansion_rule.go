package rule

import (
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
)

type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) Name() string { return "EdgeExpansionRule" }

func (r EdgeExpansionRule) Narrow(cells bits.Cells, d domain.LineDomain) (domain.LineDomain, bool) {
	if d.RunsCount() < 1 {
		return d, false
	}
	mostLeftBlack := cells.Blacks.LeftZeros()
	run := d.Run(0)
	if !run.CoversLeft(mostLeftBlack) {
		return d, false
	}
	run.MaxStart = mostLeftBlack
	return d.Narrowed(0, run)
}
