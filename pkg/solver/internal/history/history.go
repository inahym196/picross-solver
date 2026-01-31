package history

import (
	"iter"
	"slices"

	"github.com/inahym196/picross-solver/pkg/solver/internal/domain"
)

type Step struct {
	RuleName string
	Domain   domain.LineDomain
}

type History struct {
	steps []Step
}

func NewHistory() *History { return &History{} }

func (h *History) Append(steps ...Step) {
	h.steps = append(h.steps, steps...)
}

func (h *History) Merge(other History) {
	h.steps = append(h.steps, other.steps...)
}

func (h *History) IsEmpty() bool { return len(h.steps) == 0 }

func (h *History) Last() Step { return h.steps[len(h.steps)-1] }

func (h *History) All() iter.Seq2[int, Step] {
	return slices.All(h.steps)
}
