package solver

import "github.com/inahym196/picross-solver/pkg/game"

type LineRefQueue struct {
	q   []game.LineRef
	inQ map[game.LineRef]struct{}
}

func (rq *LineRefQueue) Push(ref game.LineRef) {
	if rq.inQ == nil {
		rq.inQ = make(map[game.LineRef]struct{})
	}
	if _, ok := rq.inQ[ref]; ok {
		return
	}
	rq.q = append(rq.q, ref)
	rq.inQ[ref] = struct{}{}
}

func (rq *LineRefQueue) Pop() (game.LineRef, bool) {
	if len(rq.q) == 0 {
		return game.LineRef{}, false
	}
	ref := rq.q[0]
	rq.q = rq.q[1:]
	delete(rq.inQ, ref)
	return ref, true
}
