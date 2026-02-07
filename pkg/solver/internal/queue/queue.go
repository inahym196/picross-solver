package queue

type Queue[T comparable] struct {
	q   []T
	inQ map[T]struct{}
}

func (rq *Queue[T]) Push(ref T) {
	if rq.inQ == nil {
		rq.inQ = make(map[T]struct{})
	}
	if _, ok := rq.inQ[ref]; ok {
		return
	}
	rq.q = append(rq.q, ref)
	rq.inQ[ref] = struct{}{}
}

func (rq *Queue[T]) Pop() (T, bool) {
	if len(rq.q) == 0 {
		var zero T
		return zero, false
	}
	ref := rq.q[0]
	rq.q = rq.q[1:]
	delete(rq.inQ, ref)
	return ref, true
}
