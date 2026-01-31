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

func (run RunPlacement) ForcedMask() bits.Bits {
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

const MaxRuns = 16 // uint32, 32/2 = 16

type RunPlacements struct {
	count int
	runs  [MaxRuns]RunPlacement
}

func (runs RunPlacements) Equals(other RunPlacements) bool { return runs == other }
func (runs RunPlacements) Count() int                      { return runs.count }
func (runs RunPlacements) At(i int) (RunPlacement, error) {
	if !runs.inBounds(i) {
		return RunPlacement{}, fmt.Errorf("out of range: %d", i)
	}
	return runs.runs[i], nil
}

func (runs RunPlacements) Append(run RunPlacement) error {
	if runs.count >= MaxRuns {
		return fmt.Errorf("capacity over. maxRuns: %d", MaxRuns)
	}
	runs.runs[runs.count] = run
	runs.count++
	return nil
}

func (runs RunPlacements) ForcedMask() bits.Bits {
	var mask bits.Bits
	for i := range runs.count {
		mask |= runs.runs[i].ForcedMask()
	}
	return mask
}
func (runs RunPlacements) Replaced(i int, run RunPlacement) (RunPlacements, error) {
	if !runs.inBounds(i) {
		return RunPlacements{}, fmt.Errorf("out of range")
	}
	runs.runs[i] = run
	return runs, nil
}

func (runs RunPlacements) inBounds(i int) bool { return 0 <= i && i < runs.count }
