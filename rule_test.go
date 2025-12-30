package picrosssolver

import (
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
