package domain

import (
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
)

// ValueObject
type LineDomain struct {
	lineLen int
	runs    RunPlacements
}

func NewLineDomain(lineLen int, hints []int) LineDomain {
	runs := RunPlacements{}

	var sum int
	for _, h := range hints {
		sum += h
	}
	minLength := sum + len(hints) - 1
	margin := lineLen - minLength
	start := 0
	for _, hint := range hints {
		runs.Append(RunPlacement{
			MinStart: start,
			MaxStart: start + margin,
			Len:      hint,
		})
		start += hint + 1
	}
	return LineDomain{lineLen, runs}
}

func (ld LineDomain) Len() int                        { return ld.runs.count }
func (ld LineDomain) Run(i int) (RunPlacement, error) { return ld.runs.At(i) }
func (ld LineDomain) Equals(other LineDomain) bool    { return ld.runs.Equals(other.runs) }

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

func (ld LineDomain) NarrowedMax(i int, maxStart int) LineDomain {
	run, err := ld.runs.At(i)
	if err != nil {
		panic("invalid")
	}
	if maxStart >= run.MaxStart {
		panic("invalid")
	}
	run.MaxStart = maxStart
	ld.runs, err = ld.runs.Replaced(i, run)
	if err != nil {
		panic("invalid")
	}
	return ld
}
