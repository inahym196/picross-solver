package domain

import (
	"fmt"

	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
)

type RunPlacement struct {
	StartCandidates bits.Bits
	Len             int
}

func (run RunPlacement) String() string {
	return fmt.Sprintf("{Len:%d Starts:%b}", run.Len, run.StartCandidates)
}

func (run RunPlacement) minStart() int {
	return run.StartCandidates.TailZeros()
}

func (run RunPlacement) maxStart() int {
	return bits.Bits(0).UnitSize() - run.StartCandidates.HeadZeros() - 1
}

func (run RunPlacement) CoveredMask() bits.Bits {

	end := run.minStart() + run.Len
	maxStart := run.maxStart()
	if maxStart >= end {
		return 0
	}

	var mask bits.Bits
	for i := maxStart; i < end; i++ {
		mask |= 1 << i
	}
	return mask
}

func (run RunPlacement) CoverableMask() bits.Bits {
	m := run.StartCandidates
	maxStart := run.maxStart()
	for i := range run.Len {
		m |= bits.Bits(1 << (maxStart + i))
	}
	return m
}

func (run RunPlacement) Coverable(i int) bool {
	return bits.Bits(1<<i)&run.CoverableMask() != 0
}

func (run RunPlacement) Fixed(start int) RunPlacement {
	if bits.Bits(1<<start)&run.StartCandidates == 0 {
		panic("invalid start")
	}
	run.StartCandidates = 1 << start
	return run
}

// sc: low [00000111100000000000000000000000] upper, len: 3

func (run RunPlacement) WithMaxStart(max int) (RunPlacement, bool) {
	if max < 0 {
		return run, false
	}
	next := run.StartCandidates & bits.Bits(1<<(max+1)-1)
	if next == 0 || next == run.StartCandidates {
		return run, false
	}
	run.StartCandidates = next
	return run, true
}

func (run RunPlacement) WithMinStart(min int) (RunPlacement, bool) {
	if min <= 0 {
		return run, false
	}
	next := run.StartCandidates &^ bits.Bits(1<<min-1)
	if next == 0 || next == run.StartCandidates {
		return run, false
	}
	run.StartCandidates = next
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
func (runs RunPlacements) At(i int) RunPlacement {
	if !runs.inBounds(i) {
		panic("out of range")
	}
	return runs.runs[i]
}
func (runs RunPlacements) Replaced(i int, newRun RunPlacement) (newRuns RunPlacements, changed bool) {
	if run := runs.At(i); run.Equals(newRun) {
		return runs, false
	}
	runs.runs[i] = newRun
	return runs, true
}

func (runs RunPlacements) FixedByMask(mask bits.Bits) (RunPlacements, bool) {
	cursor := 0
	changed := false

	for i := range runs.count {
		run := runs.runs[i]
		start := cursor + (mask >> cursor).TailZeros()
		if bits.Bits(1<<start)&run.StartCandidates == 0 {
			return runs, false
		}

		fixed := run.Fixed(start)
		if fixed == run {
			return runs, false
		}
		runs.runs[i] = fixed
		changed = true
		cursor = start + run.Len
	}
	return runs, changed
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

func (runs RunPlacements) IsExactFit() bool {
	for _, run := range runs.runs {
		if run.StartCandidates.OnesCount() != 1 {
			return false
		}
	}
	return true
}

func (runs RunPlacements) inBounds(i int) bool { return 0 <= i && i < runs.count }
