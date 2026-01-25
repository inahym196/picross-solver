package game_test

import (
	"reflect"
	"testing"

	. "github.com/inahym196/picross-solver/internal/testutil"
	"github.com/inahym196/picross-solver/pkg/game"
)

func TestBoard_Mark(t *testing.T) {
	tests := []struct {
		board *game.Board
		row   int
		col   int
		cell  game.Cell
		want  [][]game.Cell
	}{
		{game.NewBoard(2, 2), 0, 0, game.CellBlack, [][]game.Cell{{B, U}, {U, U}}},
		{game.NewBoard(2, 2), 1, 1, game.CellWhite, [][]game.Cell{{U, U}, {U, W}}},
	}

	for _, tt := range tests {
		tt.board.Mark(tt.row, tt.col, tt.cell)
		got := tt.board.Cells()

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}
}

func TestBoard_MarkMulti(t *testing.T) {
	b := game.NewBoard(2, 2)
	want := [][]game.Cell{{W, W}, {B, B}}

	b.Mark(0, 0, game.CellWhite)
	b.Mark(0, 1, game.CellWhite)
	b.Mark(1, 0, game.CellBlack)
	b.Mark(1, 1, game.CellBlack)
	got := b.Cells()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %+v, got: %+v", want, got)

	}
}

func TestGame_Mark(t *testing.T) {
	g, err := game.NewGame([][]int{{2}, {0}}, [][]int{{1}, {1}})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	want := [][]game.Cell{{B, U}, {U, U}}

	g.Mark(0, 0, game.CellBlack)
	got := g.Cells()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGame_LinesMark(t *testing.T) {
	g, err := game.NewGame([][]int{{2}, {0}}, [][]int{{1}, {1}})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	want := [][]game.Cell{{W, W}, {B, B}, {W, B}, {W, B}}

	g.Mark(0, 0, game.CellWhite)
	g.Mark(0, 1, game.CellWhite)
	g.Mark(1, 0, game.CellBlack)
	g.Mark(1, 1, game.CellBlack)
	got := g.Lines()

	if len(got) != 4 {
		t.Errorf("want %d, got %d", 4, len(got))
	}
	for i := range got {
		if !reflect.DeepEqual(got[i].Cells, want[i]) {
			t.Errorf("%d: want %v, got %v", i, want[i], got[i].Cells)
		}
	}

	wantCells := [][]game.Cell{{W, W}, {B, B}}
	gotCells := g.Cells()
	if !reflect.DeepEqual(gotCells, wantCells) {
		t.Errorf("want %v, got %v", wantCells, gotCells)
	}
}
