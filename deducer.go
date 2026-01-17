package picrosssolver

import (
	"fmt"
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

type deduction struct {
	ruleName string
	hints    []int
	lineRef  lineRef
	before   []game.Cell
	after    []game.Cell
}

func (deduction deduction) String() string {
	return fmt.Sprintf("%s %s %v %v -> %v", deduction.ruleName, deduction.lineRef, deduction.hints, deduction.before, deduction.after)
}

type Rule interface {
	Name() string
	Deduce(lineView) []game.Cell
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
			hogehogeRule{},
			PruneImpossibleSegmentRule{},
			FillRemainingWhiteRule{},
		},
	}
}

func (d deducer) DeduceLine(line lineView, ref lineRef) (deds []deduction) {
	current := line

	for _, rule := range d.rules {
		if current.IsFilled() {
			return deds
		}

		before := slices.Clone(current.Cells)
		updated := rule.Deduce(current)

		if updated == nil || slices.Equal(before, updated) {
			continue
		}

		deds = append(deds, deduction{
			ruleName: rule.Name(),
			hints:    current.Hints,
			lineRef:  ref,
			before:   before,
			after:    updated,
		})
		current.Cells = updated
	}
	return deds
}
