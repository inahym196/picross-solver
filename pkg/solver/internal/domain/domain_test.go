package domain_test

import (
	"testing"

	"github.com/inahym196/picross-solver/pkg/solver/internal/bits"
	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
)

func str2bit(s string) bits.Bits {
	var b bits.Bits
	for i, r := range s {
		if r == '1' {
			b |= 1 << i
		}
	}
	return b
}

func str2Cells(lineLen int, blackStr, whiteStr string) bits.Cells {
	b := str2bit(blackStr)
	w := str2bit(whiteStr)
	c, err := bits.FromMasks(lineLen, b, w)
	if err != nil {
		panic(err)
	}
	return c
}

func mustNewLineDomain(lineLen int, hints []int) domain.LineDomain {
	d, err := domain.NewLineDomain(lineLen, hints)
	if err != nil {
		panic(err)
	}
	return d
}

func TestLineDomain_Project_Trivial(t *testing.T) {
	tests := []struct {
		name   string
		domain domain.LineDomain
		want   bits.Cells
	}{
		{"zero-hint/len0", mustNewLineDomain(0, []int{}), str2Cells(0, "0", "0")},
		{"zero-hint/len1", mustNewLineDomain(1, []int{0}), str2Cells(1, "0", "1")},
		{"zero-hint/len2", mustNewLineDomain(2, []int{0}), str2Cells(2, "0", "11")},
		{"zero-hint/len30", mustNewLineDomain(30, []int{0}), str2Cells(30, "0", "111111111111111111111111111111")},
		{"full-hint/len1", mustNewLineDomain(1, []int{1}), str2Cells(1, "1", "0")},
		{"full-hint/len2", mustNewLineDomain(2, []int{2}), str2Cells(2, "11", "0")},
		{"full-hint/len30", mustNewLineDomain(30, []int{30}), str2Cells(30, "111111111111111111111111111111", "0")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.domain.Project()
			if err != nil {
				t.Fatal(err)
			}

			if !got.Equals(tt.want) {
				t.Errorf("project: want %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestLineDomain_Project_Fits(t *testing.T) {
	tests := []struct {
		name   string
		domain domain.LineDomain
		want   bits.Cells
	}{
		{"fits: 1 1", mustNewLineDomain(3, []int{1, 1}), str2Cells(3, "101", "010")},
		{"fits: 2 1", mustNewLineDomain(4, []int{2, 1}), str2Cells(4, "1101", "0010")},
		{"fits: 1 1 1", mustNewLineDomain(5, []int{1, 1, 1}), str2Cells(5, "10101", "01010")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.domain.Project()
			if err != nil {
				t.Fatal(err)
			}

			if !got.Equals(tt.want) {
				t.Errorf("project: want %+v, got %+v", tt.want, got)
			}
		})
	}
}
