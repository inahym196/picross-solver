package rule_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver"
	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
	"github.com/inahym196/picross-solver/pkg/solver/internal/rule"
)

const (
	U = game.CellUndetermined
	W = game.CellWhite
	B = game.CellBlack
)

func mustNewLineDomain(len int, hint []int) domain.LineDomain {
	d, err := domain.NewLineDomain(len, hint)
	if err != nil {
		panic(err)
	}
	return d
}

func TestAllRuleV2(t *testing.T) {
	tests := []struct {
		rule      solver.RuleV2
		cells     bits.Cells
		hint      []int
		wantCells bits.Cells
	}{
		{
			rule.EdgeExpansionRule{},
			bits.FromCells([]game.Cell{U, B, U, U, U, U}),
			[]int{3},
			bits.FromCells([]game.Cell{U, B, B, U, W, W}),
		},
		{
			rule.EdgeExpansionRule{},
			bits.FromCells([]game.Cell{U, U, U, U, B, U}),
			[]int{3},
			bits.FromCells([]game.Cell{W, W, U, B, B, U}),
		},
		{
			rule.EdgeExpansionRule{},
			bits.FromCells([]game.Cell{B, U, U, U, U, U, U, B}),
			[]int{2, 1, 2},
			bits.FromCells([]game.Cell{B, B, W, U, U, W, B, B}),
		},
		{
			rule.ExactMatchBlackRule{},
			bits.FromCells([]game.Cell{B, W, W, W, U}),
			[]int{1, 1},
			bits.FromCells([]game.Cell{B, W, W, W, B}),
		},
		{
			rule.ExactMatchBlackRule{},
			bits.FromCells([]game.Cell{W, U, U, W, U, W}),
			[]int{2, 1},
			bits.FromCells([]game.Cell{W, B, B, W, B, W}),
		},
		{
			rule.ExactMatchWhiteRule{},
			bits.FromCells([]game.Cell{U, B, B, U, B, U}),
			[]int{2, 1},
			bits.FromCells([]game.Cell{W, B, B, W, B, W}),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-case%d", tt.rule.Name(), i), func(t *testing.T) {

			d := mustNewLineDomain(tt.cells.Len, tt.hint)
			got, changed := tt.rule.Narrow(tt.cells, d)
			gotProject, err := got.Project()
			if err != nil {
				t.Fatal(err)
			}

			wantChanged := !reflect.DeepEqual(tt.cells, tt.wantCells)
			if wantChanged != changed {
				t.Errorf("want Changed: %t, got %t", wantChanged, changed)
			}

			if tt.wantCells != gotProject {
				t.Errorf("expected %v, got %v", tt.wantCells, gotProject)
			}
		})
	}
}
