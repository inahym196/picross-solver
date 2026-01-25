package deducer

import (
	"fmt"

	"github.com/inahym196/picross-solver/pkg/game"
)

type Deduction struct {
	RuleName string
	Hints    []int
	LineRef  game.LineRef
	Before   []game.Cell
	After    []game.Cell
}

func (deduction Deduction) String() string {
	return fmt.Sprintf("%s %s %v %v -> %v", deduction.RuleName, deduction.LineRef, deduction.Hints, deduction.Before, deduction.After)
}

type Rule interface {
	Name() string
	Deduce(game.Line) []game.Cell
}
