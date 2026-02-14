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

func TestAllRuleV2(t *testing.T) {
	tests := []struct {
		rule       solver.RuleV2
		cells      bits.Cells
		domainFunc func() (domain.LineDomain, error)
		wantCells  bits.Cells
	}{
		{
			rule.EdgeExpansionRule{},
			bits.FromCells([]game.Cell{U, B, U, U, U, U}),
			func() (domain.LineDomain, error) { return domain.NewLineDomain(6, []int{3}) },
			bits.FromCells([]game.Cell{U, B, B, U, W, W}),
		},
		{
			rule.ExactMatchBlackRule{},
			bits.FromCells([]game.Cell{W, U, U, W, U, W}),
			func() (domain.LineDomain, error) { return domain.NewLineDomain(6, []int{2, 1}) },
			bits.FromCells([]game.Cell{W, B, B, W, B, W}),
		},
		{
			rule.ExactMatchWhiteRule{},
			bits.FromCells([]game.Cell{U, B, B, U, B, U}),
			func() (domain.LineDomain, error) { return domain.NewLineDomain(6, []int{2, 1}) },
			bits.FromCells([]game.Cell{W, B, B, W, B, W}),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-case%d", tt.rule.Name(), i), func(t *testing.T) {
			domain, err := tt.domainFunc()
			if err != nil {
				t.Fatal(err)
			}

			got, changed := tt.rule.Narrow(tt.cells, domain)
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

func TestEdgeExpansionRule_RightNarrowTargetsLastRun(t *testing.T) {
	d, err := domain.NewLineDomain(6, []int{1, 3})
	if err != nil {
		t.Fatal(err)
	}

	cells := bits.FromCells([]game.Cell{U, U, U, U, B, U})
	got, changed := rule.EdgeExpansionRule{}.Narrow(cells, d)
	if !changed {
		t.Fatal("want changed")
	}

	run0 := got.Run(0)
	if !run0.Equals(domain.RunPlacement{MinStart: 0, MaxStart: 1, Len: 1}) {
		t.Fatalf("run0 changed unexpectedly: %+v", run0)
	}

	run1 := got.Run(1)
	if !run1.Equals(domain.RunPlacement{MinStart: 2, MaxStart: 3, Len: 3}) {
		t.Fatalf("run1 not narrowed from right: %+v", run1)
	}
}

func TestEdgeExpansionRule_NoBlackNoChange(t *testing.T) {
	d, err := domain.NewLineDomain(6, []int{2})
	if err != nil {
		t.Fatal(err)
	}

	cells := bits.FromCells([]game.Cell{U, U, U, U, U, U})
	got, changed := rule.EdgeExpansionRule{}.Narrow(cells, d)
	if changed {
		t.Fatal("want unchanged")
	}
	if !got.Equals(d) {
		t.Fatalf("domain changed unexpectedly: want %v, got %v", d, got)
	}
}
