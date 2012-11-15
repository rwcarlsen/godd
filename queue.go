
package godd

//type worker struct {
//	q *queue
//}
//
//func (w *worker) Process(set Set) {
//	out := w.q.Results
//	outcome := w.q.inp.Test(set)
//	out <- &Hist{Deltas: set, Out:outcome}
//	w.q.workers <- w
//}
//
//type queue struct {
//	workers chan *worker
//	Results chan *Hist
//	inp Input
//}
//
//func NewQueue(size int, inp Input) *queue {
//	q := &queue{workers: make(chan *worker, size), inp: inp}
//	for i := 0; i < size; i++ {
//		w := &worker{q}
//		q.workers <- w
//	}
//	return q
//}
//
//func (q *queue) Process(sets ...Set) {
//	q.Results = make(chan *Hist, len(sets))
//	go func() {
//		for _, set := range sets {
//			w := <- q.workers
//			go w.Process(set)
//		}
//	}()
//}

type queue struct {
	in chan Set
	out chan *Hist
	workers []*worker
	inp Input
}

func NewQueue(size int, inp Input) *queue {
	q := &queue{in: make(chan Set, size), out: make(chan *Hist), inp: inp}
	for i := 0; i < size; i++ {
		w := &worker{q}
		q.workers = append(q.workers, w)
		go w.Run()
	}
	return q
}

type worker struct {
	q *queue
}

func (w *worker) Run() {
	for {
		set := <- w.q.in
		outcome := w.q.inp.Test(set)
		w.q.out <- &Hist{Deltas: set, Out:outcome}
	}
}

