package picrosssolver_test

import (
	"reflect"
	"testing"

	picrosssolver "github.com/inahym196/picross-solver"
)

func TestMain(t *testing.T) {
	t.Run("row test full", func(t *testing.T) {
		rowHints := [][]int{{0}, {2}}
		colHints := [][]int{{1}, {1}}
		solve := picrosssolver.Solve(rowHints, colHints).Print()
		expected := []string{
			"__",
			"##",
		}
		if !reflect.DeepEqual(solve, expected) {
			t.Errorf("expected %v, got %v", expected, solve)
		}
	})
	t.Run("col test full", func(t *testing.T) {
		rowHints := [][]int{{1}, {1}}
		colHints := [][]int{{2}, {0}}
		solve := picrosssolver.Solve(rowHints, colHints).Print()
		expected := []string{
			"#_",
			"#_",
		}
		if !reflect.DeepEqual(solve, expected) {
			t.Errorf("expected %v, got %v", expected, solve)
		}
	})
}
