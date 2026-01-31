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
		wantFunc   func() (domain.LineDomain, bool)
	}{
		{
			rule.EdgeExpansionRule{},
			bits.FromCells([]game.Cell{U, B, U, U, U, U}),
			func() (domain.LineDomain, error) { return domain.NewLineDomain(6, []int{3}) },
			func() (domain.LineDomain, bool) {
				domain, _ := domain.NewLineDomain(6, []int{3})
				return domain.NarrowedRunMax(0, 1)
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-case%d", tt.rule.Name(), i), func(t *testing.T) {

			want, wantChanged := tt.wantFunc()
			domain, err := tt.domainFunc()
			if err != nil {
				t.Fatal(err)
			}

			got, changed := tt.rule.Narrow(tt.cells, domain)

			if wantChanged != changed {
				t.Errorf("want Changed: %T, got %T", wantChanged, changed)
			}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	}
}
