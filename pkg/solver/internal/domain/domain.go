package domain

import (
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
)

// ValueObject
type LineDomain struct {
	lineLen int
	runs    RunPlacements
}

func NewLineDomain(lineLen int, hints []int) (LineDomain, error) {
	runs := RunPlacements{}

	var sum int
	for _, h := range hints {
		sum += h
	}
	minLength := sum + len(hints) - 1
	margin := lineLen - minLength
	start := 0
	var err error = nil
	for _, hint := range hints {
		runs, err = runs.Append(RunPlacement{
			MinStart: start,
			MaxStart: start + margin,
			Len:      hint,
		})
		if err != nil {
			return LineDomain{}, err
		}
		start += hint + 1
	}
	return LineDomain{lineLen, runs}, nil
}

func (ld LineDomain) LineLen() int                   { return ld.lineLen }
func (ld LineDomain) RunsCount() int                 { return ld.runs.Count() }
func (ld LineDomain) Run(i int) (RunPlacement, bool) { return ld.runs.At(i) }
func (ld LineDomain) Equals(other LineDomain) bool {
	return ld.lineLen == other.lineLen && ld.runs.Equals(other.runs)
}

func (ld LineDomain) Project() bits.Cells {

	if count := ld.runs.Count(); count == 1 {
		run, _ := ld.runs.At(0)
		switch run.Len {
		case 0:
			return bits.NewCellsWithWhiteMasked(ld.lineLen)
		case count:
			return bits.NewCellsWithBlackMasked(ld.lineLen)
		}
	}
	return bits.NewCells(ld.lineLen).MarkedBlacks(ld.runs.ForcedMask())
}

func (ld LineDomain) NarrowedRunMax(i int, maxStart int) (LineDomain, bool) {
	run, ok := ld.runs.At(i)
	if !ok {
		return ld, false
	}
	newRun, changed := run.WithMaxStart(maxStart)
	if !changed {
		return ld, false
	}
	newRuns, ok := ld.runs.Replaced(i, newRun)
	if !ok || newRun.Equals(run) {
		return ld, false
	}
	return LineDomain{ld.lineLen, newRuns}, true
}

func (ld LineDomain) IsDeterministic() bool {
	if ld.runs.Count() == 1 {
		run, _ := ld.runs.At(0)
		if run.Len == 0 || run.Len == ld.lineLen {
			return true
		}
	}
	return ld.runs.IsExactFit()
}
