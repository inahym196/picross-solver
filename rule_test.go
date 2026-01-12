package picrosssolver

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	U = CellUndetermined
	W = CellWhite
	B = CellBlack
)

func TestSplitByWhite(t *testing.T) {
	tests := []struct {
		cells    []Cell
		expected [][]Cell
	}{
		{[]Cell{U, W, U}, [][]Cell{{U}, {U}}},
		{[]Cell{W, B, U, U, W}, [][]Cell{{B, U, U}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			got := splitByWhite(tt.cells)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestAllRule(t *testing.T) {
	tests := []struct {
		rule     Rule
		cells    []Cell
		hints    []int
		expected []Cell
	}{
		{ZeroHintRule{}, []Cell{U, U, U}, []int{0}, []Cell{W, W, W}},
		{MinimumSpacingRule{}, []Cell{U, U, U}, []int{1, 1}, []Cell{B, W, B}},
		{MinimumSpacingRule{}, []Cell{U, U, U}, []int{3}, []Cell{B, B, B}},
		{MinimumSpacingRule{}, []Cell{U, U, U, U}, []int{2, 1}, []Cell{B, B, W, B}},
		{MinimumSpacingRule{}, []Cell{W, U, U, U, U}, []int{1, 2}, []Cell{W, B, W, B, B}},
		{MinimumSpacingRule{}, []Cell{U, U, U, U, U}, []int{1, 1, 1}, []Cell{B, W, B, W, B}},
		{MinimumSpacingRule{}, []Cell{U, U, U, U, U, U}, []int{1, 2, 1}, []Cell{B, W, B, B, W, B}},
		{OverlapFillRule{}, []Cell{U, U, U}, []int{2}, []Cell{U, B, U}},
		{OverlapFillRule{}, []Cell{U, U, U, U}, []int{3}, []Cell{U, B, B, U}},
		{OverlapFillRule{}, []Cell{U, U, U, U, U}, []int{2, 1}, []Cell{U, B, U, U, U}},
		{OverlapFillRule{}, []Cell{W, U, W, U, U}, []int{1, 2}, []Cell{W, B, W, B, B}},
		{OverlapFillRule{}, []Cell{W, U, W, U, U, U}, []int{1, 2}, []Cell{W, B, W, U, B, U}},
		{OverlapFillRule{}, []Cell{U, U, U, U, U, U}, []int{2, 2}, []Cell{U, B, U, U, B, U}},
		{OverlapFillRule{}, []Cell{U, U, U, W, U, U}, []int{1, 2}, []Cell{U, U, U, W, B, B}},
		{OverlapFillRule{}, []Cell{B, B, W, U, U, U, B, B, W, U, U, U, W, U, U}, []int{2, 4, 1, 1}, []Cell{B, B, W, U, B, B, B, B, W, U, U, U, W, U, U}},
		{OverlapExpansionRule{}, []Cell{U, U, U, U, U}, []int{1, 1}, nil},
		{OverlapExpansionRule{}, []Cell{U, B, U, U, U, U}, []int{3}, []Cell{U, B, B, U, U, U}},
		{OverlapExpansionRule{}, []Cell{U, U, U, U, B, U}, []int{3}, []Cell{U, U, U, B, B, U}},
		{OverlapExpansionRule{}, []Cell{W, U, B, U, U, U, U}, []int{3}, []Cell{W, U, B, B, U, U, U}},
		{EdgeExpansionRule{}, []Cell{B, U, U}, []int{2}, []Cell{B, B, W}},
		{EdgeExpansionRule{}, []Cell{U, U, B}, []int{2}, []Cell{W, B, B}},
		{EdgeExpansionRule{}, []Cell{W, B, U, U}, []int{2}, []Cell{W, B, B, W}},
		{EdgeExpansionRule{}, []Cell{U, U, B, W}, []int{2}, []Cell{W, B, B, W}},
		{EdgeExpansionRule{}, []Cell{W, W, B, U, U, U}, []int{3}, []Cell{W, W, B, B, B, W}},
		{EdgeExpansionRule{}, []Cell{U, U, U, B, W, W}, []int{3}, []Cell{W, B, B, B, W, W}},
		{BlockSatisfiedRule{}, []Cell{U}, []int{1, 1}, nil},
		{BlockSatisfiedRule{}, []Cell{U, B, U}, []int{1}, []Cell{W, B, W}},
		{BlockSatisfiedRule{}, []Cell{U, U, B, B, U, U}, []int{2}, []Cell{U, W, B, B, W, U}},
		{BlockSatisfiedRule{}, []Cell{B, U, U, B, B, U}, []int{1, 2}, []Cell{B, U, W, B, B, W}},
		{PruneImpossibleSegmentRule{}, []Cell{U, W, U}, []int{1}, nil},
		{PruneImpossibleSegmentRule{}, []Cell{U, W, U, U}, []int{2}, []Cell{W, W, U, U}},
		{PruneImpossibleSegmentRule{}, []Cell{U, W, U, W, U, U}, []int{2, 3, 4, 5}, []Cell{W, W, W, W, U, U}},
		{FillRemainingWhiteRule{}, []Cell{U, B, U}, []int{1}, []Cell{W, B, W}},
		{FillRemainingWhiteRule{}, []Cell{B, U, B, B}, []int{1, 2}, []Cell{B, W, B, B}},
		{FillRemainingWhiteRule{}, []Cell{U, B, B, W, U, B}, []int{2, 1}, []Cell{W, B, B, W, W, B}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-case%d", tt.rule.Name(), i), func(t *testing.T) {
			line := lineView{tt.cells, tt.hints}

			got := tt.rule.Deduce(line)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
