package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/inahym196/picross-solver/pkg/game"
	"github.com/inahym196/picross-solver/pkg/solver"
)

func parseHints(s string) [][]int {
	fields := strings.Fields(s)
	hints := make([][]int, 0, len(fields))
	for _, f := range fields {
		parts := strings.Split(f, "-")
		row := make([]int, 0, len(parts))
		for _, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				panic(err)
			}
			row = append(row, n)
		}
		hints = append(hints, row)
	}
	return hints
}

func main() {

	rowHints := parseHints(`
	6 2-4-5 1-1-2-1-5-4 1-1-1-2-3-1-2 1-1-2-5-1
	5-1-2-2-1 1-2-6-6 2-1-3-4-2 3-1-3-2 4-1-1-2-3
	1-3-1-1-2 1-6-1-3-2-2 2-3-1-2-2-2 1-1-1-1-4-3 1-3-1-4-5
	1-1-2-1-3-6 1-1-2-2-4-3 1-4-1-6 2-2-2-3-4 1-3-1-5-4
	2-2-3-3-3 2-4-2-3-3 2-3-2-3-2 6-6 5-5`)
	colHints := parseHints(`
	4-7 1-3-2-3 3-1-2-1-1-2 1-1-2-3-2-2 1-2-1-2-2-2-1-1-2
	3-1-1-9-1-1 2-1-1-1-1-9 2-1-1-1-2-5 2-1-2-4 3-6-2
	2-3-2-1 1-2-1-2 2-1-3 3-2-1-1-2-3 4-3-4-1-2
	1-2-3-8-4 1-1-2-7-1 1-1-2-6-1 4-3-7 1-2-3-1-1-6
	2-1-4-1-4 1-2-1-3-2-2 1-6-3-7 2-2-2-3-5-1 1-1-1-3-4`)
	game, _ := game.NewGame(rowHints, colHints)
	solver := solver.NewSolver()

	n, deds := solver.ApplyMany(game)
	fmt.Println("logs: ")
	for i, log := range deds {
		fmt.Printf("%2d: %v\n", i, log)
	}
	fmt.Printf("applied x%d\n", n)

	for i, s := range game.Board().Print() {
		fmt.Printf("%2d: %s\n", i, s)
	}

}
