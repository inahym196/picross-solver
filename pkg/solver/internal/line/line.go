package line

import (
	"iter"
	"slices"

	"github.com/inahym196/picross-solver/pkg/game"
)

type Line struct {
	Cells []game.Cell
	Hints []int
}

func (l Line) IsFilled() bool {
	return slices.Index(l.Cells, game.CellUndetermined) == -1
}

type Span struct {
	MinStart int
	MaxStart int
	Length   int
}

type LineConstraint struct {
	spans []Span
}

func NewLineConstraint(length int, hints []int) LineConstraint {
	spans := make([]Span, 0, len(hints))

	var sum int
	for _, h := range hints {
		sum += h
	}
	minTotal := sum + len(hints) - 1
	margin := length - minTotal
	start := 0
	for _, hint := range hints {
		spans = append(spans, Span{
			MinStart: start,
			MaxStart: start + margin,
			Length:   hint,
		})
		start += hint + 1
	}
	return LineConstraint{spans}
}

func (lc LineConstraint) All() iter.Seq2[int, Span] { return slices.All(lc.spans) }
