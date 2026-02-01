package solver_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver"
	"github.com/inahym196/picross-solver/pkg/solver/internal/history"
)

func TestV2E2E(t *testing.T) {
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
		}, {
			rowHints: ParseHints("2-3-1-2-3 1-2-4-1 1-2-5 3-2-2-1 1-1-2-1-1 4-1-1-2 5-1-1-3 5-1-1-3 2-1-1-1-1-1 1-1-1-1-1-1 2-1-3 1-8-1 0 1-1-1-1-1-1 2-2"),
			colHints: ParseHints("1-8-2 1-1-4-1-1 1-3-1 2-4-3-1 1-1-3-1 4-1 1-2-3-1 1-5-1 2-1 3-6-1 6-2 3-3-1 1-1-2 2-4-1-1 1-1-5-2"),
			expected: []string{
				"##_###_#_##_###",
				"___#_##_####_#_",
				"#___##__#####__",
				"###__##___##__#",
				"#__#__##__#__#_",
				"####___#__#__##",
				"#####__#_#__###",
				"#####__#_#__###",
				"##__#__#_#_#__#",
				"#__#__#__#_#__#",
				"__##__#__###___",
				"_#_########__#_",
				"_______________",
				"#__#__#__#_#__#",
				"##___________##",
			},
		},
	}
	solver := solver.NewSolverV2()
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			game, _ := game.NewGame(tt.rowHints, tt.colHints)
			h := history.NewHistory()

			n := solver.ApplyMany(game, h)
			t.Logf("applied x%d\n", n)
			boardStrings := game.Board().Print()
			if !reflect.DeepEqual(boardStrings, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, boardStrings)
				if h.IsEmpty() {
					t.Log("history: nil")
					t.SkipNow()
				}
				t.Log("logs: ")
				for _, h := range h.All() {
					t.Logf("  %+v\n", h)
				}
			}
		})
	}
}

func BenchmarkV2E2E(b *testing.B) {

	rowHints := ParseHints("2-3-1-2-3 1-2-4-1 1-2-5 3-2-2-1 1-1-2-1-1 4-1-1-2 5-1-1-3 5-1-1-3 2-1-1-1-1-1 1-1-1-1-1-1 2-1-3 1-8-1 0 1-1-1-1-1-1 2-2")
	colHints := ParseHints("1-8-2 1-1-4-1-1 1-3-1 2-4-3-1 1-1-3-1 4-1 1-2-3-1 1-5-1 2-1 3-6-1 6-2 3-3-1 1-1-2 2-4-1-1 1-1-5-2")
	solver := solver.NewSolverV2()
	game, _ := game.NewGame(rowHints, colHints)
	h := history.NewHistory()

	for b.Loop() {
		solver.ApplyMany(game, h)
	}

}
