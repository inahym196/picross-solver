package picrosssolver_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	picrosssolver "github.com/inahym196/picross-solver"
)

func ParseHints(s string) [][]int {
	fields := strings.Fields(s)
	hints := make([][]int, 0, len(fields))
	for _, f := range fields {
		parts := strings.Split(f, "-")
		row := make([]int, 0, len(parts))
		for _, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				panic(err)
			}
			row = append(row, n)
		}
		hints = append(hints, row)
	}
	return hints
}

func TestParseHints(t *testing.T) {
	s := `0-1 2`
	got := ParseHints(s)
	expected := [][]int{{0, 1}, {2}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestE2E(t *testing.T) {
	tests := []struct {
		rowHints [][]int
		colHints [][]int
		expected []string
	}{
		{
			rowHints: ParseHints("0 2"),
			colHints: ParseHints("1 1"),
			expected: []string{
				"__",
				"##",
			},
		},
		{
			rowHints: ParseHints("1 1"),
			colHints: ParseHints("2 0"),
			expected: []string{
				"#_",
				"#_",
			},
		},
		{
			rowHints: ParseHints("1-1-1 1-1-1 5 5 5"),
			colHints: ParseHints("5 3 5 3 5"),
			expected: []string{
				"#_#_#",
				"#_#_#",
				"#####",
				"#####",
				"#####",
			},
		},
		{
			rowHints: ParseHints("5 1-1 1-1 1-1 1-2"),
			colHints: ParseHints("1 5 1 5 1-1"),
			expected: []string{
				"#####",
				"_#_#_",
				"_#_#_",
				"_#_#_",
				"_#_##",
			},
		},
		{
			rowHints: ParseHints("1 1 5 1 1"),
			colHints: ParseHints("1 3 1-1-1 1 1"),
			expected: []string{
				"__#__",
				"_#___",
				"#####",
				"_#___",
				"__#__",
			},
		},
	}
	solver := picrosssolver.NewSolver()
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			game, _ := picrosssolver.NewGame(tt.rowHints, tt.colHints)

			solve, n, logs := solver.ApplyMany(*game)
			t.Logf("applied x%d\n", n)

			if !reflect.DeepEqual(solve.Print(), tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, solve.Print())
				t.Log("logs: ")
				for _, log := range logs {
					t.Logf("  %+v\n", log)
				}
			}

		})
	}
}
