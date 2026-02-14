package rule

import (
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
)

type EdgeExpansionRule struct{}

func (r EdgeExpansionRule) Name() string { return "EdgeExpansionRule" }

func (r EdgeExpansionRule) leftNarrow(cells bits.Cells, d domain.LineDomain) (domain.LineDomain, bool) {
	if cells.Blacks == 0 {
		return d, false
	}
	mostLeftBlack := cells.Blacks.LeftZeros()
	run := d.Run(0)
	if !run.CoversLeft(mostLeftBlack) {
		return d, false
	}
	run, changed := run.WithMaxStart(mostLeftBlack)
	if !changed {
		return d, false
	}
	return d.Narrowed(0, run)
}

func (r EdgeExpansionRule) rightNarrow(cells bits.Cells, d domain.LineDomain) (domain.LineDomain, bool) {
	if cells.Blacks == 0 {
		return d, false
	}
	mostRightBlack := 31 - cells.Blacks.RightZeros()
	last := d.RunsCount() - 1
	run := d.Run(last)
	if !run.CoversRight(mostRightBlack) {
		return d, false
	}
	minStart := mostRightBlack - run.Len + 1
	run, changed := run.WithMinStart(minStart)
	if !changed {
		return d, false
	}
	return d.Narrowed(last, run)
}

func (r EdgeExpansionRule) Narrow(cells bits.Cells, d domain.LineDomain) (domain.LineDomain, bool) {
	if d.RunsCount() < 1 {
		return d, false
	}
	var leftChanged, rightChanged bool
	d, leftChanged = r.leftNarrow(cells, d)
	d, rightChanged = r.rightNarrow(cells, d)
	return d, leftChanged || rightChanged
}
