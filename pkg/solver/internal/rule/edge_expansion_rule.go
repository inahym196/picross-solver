package rule

import (
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
)

type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) Name() string { return "EdgeExpansionRule" }

func (r EdgeExpansionRule) Narrow(cells bits.Cells, domain domain.LineDomain) (domain.LineDomain, bool) {
	if domain.RunsCount() < 1 {
		return domain, false
	}
	mostLeftBlack, found := cells.LeftMostBlackNotWhite()
	if !found {
		return domain, false
	}
	run, _ := domain.Run(0)
	if !run.CoversLeft(mostLeftBlack) {
		return domain, false
	}
	return domain.NarrowedRunMax(0, mostLeftBlack)
}
