package deducer

import (
	"fmt"
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/accessor"
	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
	"github.com/inahym196/picross-solver/pkg/solver/internal/rules"
)

type Deduction struct {
	RuleName string
	Hints    []int
	LineRef  accessor.LineRef
	Before   []game.Cell
	After    []game.Cell
}

func (deduction Deduction) String() string {
	return fmt.Sprintf("%s %s %v %v -> %v", deduction.RuleName, deduction.LineRef, deduction.Hints, deduction.Before, deduction.After)
}

type Rule interface {
	Name() string
	Deduce(line.Line) []game.Cell
}

type Deducer struct {
	rules []Rule
}

func NewDeducer() Deducer {
	return Deducer{
		[]Rule{
			rules.ZeroHintRule{},
			rules.MinimumSpacingRule{},
			rules.OverlapFillRule{},
			rules.OverlapExpansionRule{},
			rules.EdgeExpansionRule{},
			rules.BlockSatisfiedRule{},
			rules.HogeHogeRule{},
			rules.PruneImpossibleSegmentRule{},
			rules.FillRemainingWhiteRule{},
		},
	}
}

func (d Deducer) DeduceLine(line line.Line, ref accessor.LineRef) (deds []Deduction) {
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

		deds = append(deds, Deduction{
			RuleName: rule.Name(),
			Hints:    current.Hints,
			LineRef:  ref,
			Before:   before,
			After:    updated,
		})
		current.Cells = updated
	}
	return deds
}
