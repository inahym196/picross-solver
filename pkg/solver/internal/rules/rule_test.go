package rules_test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	. "github.com/inahym196/picross-solver/internal/testutil"
	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver/internal/deducer"
	"github.com/inahym196/picross-solver/pkg/solver/internal/rules"
)

func assertRuleIsPure(t *testing.T, r deducer.Rule, line game.Line) {
	t.Helper()

	origCells := slices.Clone(line.Cells)
	origHints := slices.Clone(line.Hints)

	r.Deduce(line)

	if !slices.Equal(origCells, line.Cells) {
		t.Fatalf("%s mutated Cells", r.Name())
	}
	if !slices.Equal(origHints, line.Hints) {
		t.Fatalf("%s mutated Hints", r.Name())
	}
}

func TestSplitByWhite(t *testing.T) {
	tests := []struct {
		cells    []game.Cell
		expected [][]game.Cell
	}{
		{[]game.Cell{U, W, U}, [][]game.Cell{{U}, {U}}},
		{[]game.Cell{W, B, U, U, W}, [][]game.Cell{{B, U, U}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			got := rules.SplitByWhite(tt.cells)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestAllRule(t *testing.T) {
	tests := []struct {
		rule     deducer.Rule
		cells    []game.Cell
		hints    []int
		expected []game.Cell
	}{
		{rules.ZeroHintRule{}, []game.Cell{U, U, U}, []int{0}, []game.Cell{W, W, W}},
		{rules.MinimumSpacingRule{}, []game.Cell{U, U, U}, []int{1, 1}, []game.Cell{B, W, B}},
		{rules.MinimumSpacingRule{}, []game.Cell{U, U, U}, []int{3}, []game.Cell{B, B, B}},
		{rules.MinimumSpacingRule{}, []game.Cell{U, U, U, U}, []int{2, 1}, []game.Cell{B, B, W, B}},
		{rules.MinimumSpacingRule{}, []game.Cell{W, U, U, U, U}, []int{1, 2}, []game.Cell{W, B, W, B, B}},
		{rules.MinimumSpacingRule{}, []game.Cell{U, U, U, U, U}, []int{1, 1, 1}, []game.Cell{B, W, B, W, B}},
		{rules.MinimumSpacingRule{}, []game.Cell{U, U, U, U, U, U}, []int{1, 2, 1}, []game.Cell{B, W, B, B, W, B}},
		{rules.OverlapFillRule{}, []game.Cell{U, U, U}, []int{2}, []game.Cell{U, B, U}},
		{rules.OverlapFillRule{}, []game.Cell{U, U, U, U}, []int{3}, []game.Cell{U, B, B, U}},
		{rules.OverlapFillRule{}, []game.Cell{U, U, U, U, U}, []int{2, 1}, []game.Cell{U, B, U, U, U}},
		{rules.OverlapFillRule{}, []game.Cell{W, U, W, U, U}, []int{1, 2}, []game.Cell{W, B, W, B, B}},
		{rules.OverlapFillRule{}, []game.Cell{W, U, W, U, U, U}, []int{1, 2}, []game.Cell{W, B, W, U, B, U}},
		{rules.OverlapFillRule{}, []game.Cell{U, U, U, U, U, U}, []int{2, 2}, []game.Cell{U, B, U, U, B, U}},
		{rules.OverlapFillRule{}, []game.Cell{U, U, U, W, U, U}, []int{1, 2}, []game.Cell{U, U, U, W, B, B}},
		{rules.OverlapFillRule{}, []game.Cell{B, B, W, U, U, U, B, B, W, U, U, U, W, U, U}, []int{2, 4, 1, 1}, []game.Cell{B, B, W, U, B, B, B, B, W, U, U, U, W, U, U}},
		{rules.OverlapExpansionRule{}, []game.Cell{U, U, U, U, U}, []int{1, 1}, nil},
		{rules.OverlapExpansionRule{}, []game.Cell{U, B, U, U, U, U}, []int{3}, []game.Cell{U, B, B, U, U, U}},
		{rules.OverlapExpansionRule{}, []game.Cell{U, U, U, U, B, U}, []int{3}, []game.Cell{U, U, U, B, B, U}},
		{rules.OverlapExpansionRule{}, []game.Cell{W, U, B, U, U, U, U}, []int{3}, []game.Cell{W, U, B, B, U, U, U}},
		{rules.EdgeExpansionRule{}, []game.Cell{B, U, U}, []int{2}, []game.Cell{B, B, W}},
		{rules.EdgeExpansionRule{}, []game.Cell{U, U, B}, []int{2}, []game.Cell{W, B, B}},
		{rules.EdgeExpansionRule{}, []game.Cell{W, B, U, U}, []int{2}, []game.Cell{W, B, B, W}},
		{rules.EdgeExpansionRule{}, []game.Cell{U, U, B, W}, []int{2}, []game.Cell{W, B, B, W}},
		{rules.EdgeExpansionRule{}, []game.Cell{W, W, B, U, U, U}, []int{3}, []game.Cell{W, W, B, B, B, W}},
		{rules.EdgeExpansionRule{}, []game.Cell{U, U, U, B, W, W}, []int{3}, []game.Cell{W, B, B, B, W, W}},
		{rules.MaxHintBlockBoundaryRule{}, []game.Cell{U}, []int{1, 1}, nil},
		{rules.MaxHintBlockBoundaryRule{}, []game.Cell{U, B, U}, []int{1}, []game.Cell{W, B, W}},
		{rules.MaxHintBlockBoundaryRule{}, []game.Cell{U, U, B, B, U, U}, []int{2}, []game.Cell{U, W, B, B, W, U}},
		{rules.MaxHintBlockBoundaryRule{}, []game.Cell{B, U, U, B, B, U}, []int{1, 2}, []game.Cell{B, U, W, B, B, W}},
		{rules.HogeHogeRule{}, []game.Cell{U, U, U, B, U}, []int{1, 2}, nil},
		{rules.HogeHogeRule{}, []game.Cell{B, U, U, B, U}, []int{1, 1}, []game.Cell{B, W, W, B, W}},
		{rules.HogeHogeRule{}, []game.Cell{U, U, U, U, B, B}, []int{1, 2}, []game.Cell{U, U, U, W, B, B}},
		{rules.PruneImpossibleSegmentRule{}, []game.Cell{U, W, U}, []int{1}, nil},
		{rules.PruneImpossibleSegmentRule{}, []game.Cell{U, W, U, U}, []int{2}, []game.Cell{W, W, U, U}},
		{rules.PruneImpossibleSegmentRule{}, []game.Cell{U, W, U, W, U, U}, []int{2, 3, 4, 5}, []game.Cell{W, W, W, W, U, U}},
		{rules.FillRemainingWhiteRule{}, []game.Cell{U, B, U}, []int{1}, []game.Cell{W, B, W}},
		{rules.FillRemainingWhiteRule{}, []game.Cell{B, U, B, B}, []int{1, 2}, []game.Cell{B, W, B, B}},
		{rules.FillRemainingWhiteRule{}, []game.Cell{U, B, B, W, U, B}, []int{2, 1}, []game.Cell{W, B, B, W, W, B}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-case%d", tt.rule.Name(), i), func(t *testing.T) {
			line := game.Line{Cells: tt.cells, Hints: tt.hints}
			assertRuleIsPure(t, tt.rule, line)

			got := tt.rule.Deduce(line)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
