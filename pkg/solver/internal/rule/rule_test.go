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
