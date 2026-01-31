package domain

import "testing"

func TestNewLineDomain(t *testing.T) {
	want := []RunPlacement{{0, 3, 2}, {3, 6, 1}}

	d, _ := NewLineDomain(7, []int{2, 1})

	if d.LineLen() != 7 {
		t.Fatalf("lineLen: want 7, got %d", d.lineLen)
	}

	if d.RunsCount() != len(want) {
		t.Fatalf("runsCount: want %d, got %d", len(want), d.runs.Count())
	}

	for i, wantRun := range want {
		run, ok := d.runs.At(i)
		if !ok {
			t.Fatalf("out of range run[%d]", i)
		}
		if run != wantRun {
			t.Errorf("run: want %+v, got %+v", wantRun, run)
		}
	}
}
