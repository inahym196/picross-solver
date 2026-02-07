package domain

import (
	"fmt"

	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
)

type RunPlacement struct {
	MinStart int
	MaxStart int
	Len      int
}

func (runs RunPlacement) String() string {
	return fmt.Sprintf("{Len:%d Start:%d-%d}", runs.Len, runs.MinStart, runs.MaxStart)
}

func (run RunPlacement) CoveredMask() bits.Bits {
	end := run.MinStart + run.Len
	if run.MaxStart >= end {
		return 0
	}

	var mask bits.Bits
	for i := run.MaxStart; i < end; i++ {
		mask |= 1 << i
	}
	return mask
}

func (run RunPlacement) CoverableMask() bits.Bits {
	var m bits.Bits
	for i := run.MinStart; i < run.MaxStart+run.Len; i++ {
		m |= bits.Bits(1 << i)
	}
	return m
}

func (run RunPlacement) CoversLeft(i int) bool {
	return run.MinStart <= i && i < run.MinStart+run.Len
}

func (run RunPlacement) Fixed(start int) RunPlacement {
	if !(run.MinStart <= start && start <= run.MaxStart) {
		panic("invalid start")
	}
	return RunPlacement{
		MinStart: start,
		MaxStart: start,
		Len:      run.Len,
	}
}

func (run RunPlacement) WithMaxStart(max int) (RunPlacement, bool) {
	if max < run.MinStart {
		return run, false
	}
	if max >= run.MaxStart {
		return run, false
	}
	run.MaxStart = max
	return run, true
}

func (run RunPlacement) Equals(other RunPlacement) bool { return run == other }

const MaxRuns = 16 // uint32, 32/2 = 16

type RunPlacements struct {
	count int
	runs  [MaxRuns]RunPlacement
}

func (runs RunPlacements) String() string {
	return fmt.Sprintf("%+v", runs.runs[:runs.count])
}
func (runs RunPlacements) Equals(other RunPlacements) bool { return runs == other }
func (runs RunPlacements) Count() int                      { return runs.count }
func (runs RunPlacements) At(i int) (RunPlacement, bool) {
	if !runs.inBounds(i) {
		return RunPlacement{}, false
	}
	return runs.runs[i], true
}

func (runs RunPlacements) Append(run RunPlacement) (RunPlacements, error) {
	if runs.count >= MaxRuns {
		return runs, fmt.Errorf("capacity over. maxRuns: %d", MaxRuns)
	}
	runs.runs[runs.count] = run
	runs.count++
	return runs, nil
}

func (runs RunPlacements) CoveredMask() bits.Bits {
	var mask bits.Bits
	for i := range runs.count {
		mask |= runs.runs[i].CoveredMask()
	}
	return mask
}

func (runs RunPlacements) CoverableMask() bits.Bits {
	var m bits.Bits
	for _, run := range runs.runs {
		m |= run.CoverableMask()
	}
	return m
}

func (runs RunPlacements) UnCoverableMask(lineLen int) bits.Bits {
	return bits.Bits(1<<lineLen-1) &^ runs.CoverableMask()
}

func (runs RunPlacements) Replaced(i int, run RunPlacement) (RunPlacements, bool) {
	if !runs.inBounds(i) {
		return runs, false
	}
	runs.runs[i] = run
	return runs, true
}

func (runs RunPlacements) IsExactFit() bool {
	for _, run := range runs.runs {
		if run.MinStart != run.MaxStart {
			return false
		}
	}
	return true
}

func (runs RunPlacements) inBounds(i int) bool { return 0 <= i && i < runs.count }
