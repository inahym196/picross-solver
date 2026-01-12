package picrosssolver

import "fmt"

type deduction struct {
	ruleName string
	hints    []int
	lineRef  lineRef
	before   []Cell
	after    []Cell
}

func (deduction deduction) String() string {
	return fmt.Sprintf("%s %s %v %v -> %v", deduction.ruleName, deduction.lineRef, deduction.hints, deduction.before, deduction.after)
}

type deducer struct {
	rules []Rule
}

func newDeducer() deducer {
	return deducer{
		[]Rule{
			ZeroHintRule{},
			MinimumSpacingRule{},
			OverlapFillRule{},
			OverlapExpansionRule{},
			EdgeExpansionRule{},
			BlockSatisfiedRule{},
			PruneImpossibleSegmentRule{},
			FillRemainingWhiteRule{},
		},
	}
}
