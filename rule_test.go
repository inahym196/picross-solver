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
