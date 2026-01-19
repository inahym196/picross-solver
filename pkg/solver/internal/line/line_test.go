package line_test

import (
	"reflect"
	"testing"

	"github.com/inahym196/picross-solver/pkg/solver/internal/line"
)

func TestNewLineConstraint(t *testing.T) {
	hints := []int{2, 1}
	expected := []line.Span{
		{MinStart: 0, MaxStart: 1, Length: 2},
		{MinStart: 3, MaxStart: 4, Length: 1},
	}

	lc := line.NewLineConstraint(5, hints)

	for i, span := range lc.All() {
		if !reflect.DeepEqual(expected[i], span) {
			t.Errorf("expected %v, got %v", expected[i], span)
		}
	}
}
