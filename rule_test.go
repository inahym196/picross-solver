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

func newHintedCells(length int, hints []int) HintedCells {
	return NewHintedCells(filledCells(length, U), hints)
}

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

func TestExtractMatchRule(t *testing.T) {
	hc := newHintedCells(3, []int{3})
	expected := []Cell{B, B, B}

	got := ExtractMatchRule{}.Deduce(hc)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestZeroHintRule(t *testing.T) {
	hc := newHintedCells(3, []int{0})
	expected := []Cell{W, W, W}

	got := ZeroHintRule{}.Deduce(hc)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestMinimumSpacingRule(t *testing.T) {
	tests := []struct {
		length   int
		hints    []int
		expected []Cell
	}{
		{3, []int{1, 1}, []Cell{B, W, B}},
		{4, []int{2, 1}, []Cell{B, B, W, B}},
		{5, []int{1, 1, 1}, []Cell{B, W, B, W, B}},
		{6, []int{1, 2, 1}, []Cell{B, W, B, B, W, B}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			hc := newHintedCells(tt.length, tt.hints)

			got := MinimumSpacingRule{}.Deduce(hc)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestOverlapFillRule(t *testing.T) {
	tests := []struct {
		length   int
		hints    []int
		expected []Cell
	}{
		{3, []int{2}, []Cell{U, B, U}},
		{4, []int{3}, []Cell{U, B, B, U}},
		{5, []int{2, 1}, []Cell{U, B, U, U, U}},
		{6, []int{2, 2}, []Cell{U, B, U, U, B, U}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			hc := newHintedCells(tt.length, tt.hints)

			got := OverlapFillRule{}.Deduce(hc)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestOverlapExpansionRule(t *testing.T) {
	tests := []struct {
		cells    []Cell
		hints    []int
		expected []Cell
	}{
		{[]Cell{U, B, U, U, U, U}, []int{3}, []Cell{U, B, B, U, U, U}},
		{[]Cell{U, U, U, U, B, U}, []int{3}, []Cell{U, U, U, B, B, U}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			hc := NewHintedCells(tt.cells, tt.hints)

			got := OverlapExpansionRule{}.Deduce(hc)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestEdgeExpansionRule(t *testing.T) {
	tests := []struct {
		cells    []Cell
		hints    []int
		expected []Cell
	}{
		{[]Cell{B, U, U}, []int{2}, []Cell{B, B, U}},
		{[]Cell{U, U, B}, []int{2}, []Cell{U, B, B}},
		{[]Cell{W, B, U, U}, []int{2}, []Cell{W, B, B, U}},
		{[]Cell{U, U, B, W}, []int{2}, []Cell{U, B, B, W}},
		{[]Cell{W, W, B, U, U, U}, []int{3}, []Cell{W, W, B, B, B, U}},
		{[]Cell{U, U, U, B, W, W}, []int{3}, []Cell{U, B, B, B, W, W}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			hc := NewHintedCells(tt.cells, tt.hints)

			got := EdgeExpansionRule{}.Deduce(hc)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestBlockSatisfiedRule(t *testing.T) {
	tests := []struct {
		cells    []Cell
		hints    []int
		expected []Cell
	}{
		{[]Cell{U}, []int{1, 1}, nil},
		{[]Cell{U, B, U}, []int{1}, []Cell{W, B, W}},
		{[]Cell{U, B, B, U, U}, []int{2}, []Cell{W, B, B, W, U}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			hc := NewHintedCells(tt.cells, tt.hints)

			got := BlockSatisfiedRule{}.Deduce(hc)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestPruneImpossibleSegmentRule(t *testing.T) {
	tests := []struct {
		cells    []Cell
		hints    []int
		expected []Cell
	}{
		{[]Cell{U, W, U}, []int{1}, nil},
		{[]Cell{U, W, U, U}, []int{2}, []Cell{W, W, U, U}},
		{[]Cell{U, W, U, W, U, U}, []int{2, 3, 4, 5}, []Cell{W, W, W, W, U, U}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			hc := NewHintedCells(tt.cells, tt.hints)

			got := PruneImpossibleSegmentRule{}.Deduce(hc)

			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
