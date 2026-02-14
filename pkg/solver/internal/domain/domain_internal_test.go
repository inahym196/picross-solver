package domain

import (
	"testing"

	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
)

func starts(candidates ...int) bits.Bits {
	var b bits.Bits
	for _, c := range candidates {
		b |= bits.Bits(1 << c)
	}
	return b
}

func TestNewLineDomain(t *testing.T) {
	tests := []struct {
		name       string
		lineLen    int
		hints      []int
		wantRuns   []RunPlacement
		wantLineLn int
	}{
		{
			name:       "lineLen=7 hints=2,1",
			lineLen:    7,
			hints:      []int{2, 1},
			wantLineLn: 7,
			wantRuns: []RunPlacement{
				{StartCandidates: starts(0, 1, 2, 3), Len: 2},
				{StartCandidates: starts(3, 4, 5, 6), Len: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewLineDomain(tt.lineLen, tt.hints)
			if err != nil {
				t.Fatal(err)
			}

			if d.LineLen() != tt.wantLineLn {
				t.Fatalf("lineLen: want %d, got %d", tt.wantLineLn, d.LineLen())
			}
			if d.RunsCount() != len(tt.wantRuns) {
				t.Fatalf("runsCount: want %d, got %d", len(tt.wantRuns), d.RunsCount())
			}

			for i, wantRun := range tt.wantRuns {
				if got := d.runs.At(i); got != wantRun {
					t.Errorf("run[%d]: want %+v, got %+v", i, wantRun, got)
				}
			}
		})
	}
}
