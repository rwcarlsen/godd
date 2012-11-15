package godd

import (
	"errors"
	"runtime"
)

type setHash string

type setCache map[setHash]bool

type Set []int

func (s Set) hash() setHash {
	h := make([]byte, len(s))
	for i, v := range s {
		h[i] = byte(v)
	}
	return setHash(h)
}

type Input interface {
	Test(Set) Outcome
	Len() int
}

type Hist struct {
	Deltas Set
	Out    Outcome
}

const (
	Passed Outcome = iota
	Failed
	Undetermined
)

type Outcome int

func (out Outcome) String() string {
	switch out {
	case Passed:
		return "PASS"
	case Failed:
		return "FAIL"
	case Undetermined:
		return "UNDET"
	}
	panic("invalid outcome")
}

type Config int

const (
	CkeepHist Config = 1 << iota
	CcacheTests
	Cconcurrent
)

type Run struct {
	Inp        Input
	Minimal    Set
	Hists      []*Hist
	tested     setCache
	cacheTests bool
	keepHist   bool
	concurrent bool
	iterations int
	jobs  chan Set
	results chan *Hist
}

func MinFail(inp Input, config Config) (*Run, error) {
	r := &Run{Inp: inp}
	r.tested = make(setCache)
	initialSet := IntRange(inp.Len())

	if inp.Test(initialSet) != Failed {
		return nil, errors.New("godd: Test passes with all deltas applied")
	}

	// get specs from config mask
	r.cacheTests = CcacheTests&config != 0
	r.keepHist = CkeepHist&config != 0
	if r.concurrent = Cconcurrent&config != 0; r.concurrent {
		r.jobs = make(chan Set, runtime.NumCPU())
		r.results = make(chan *Hist, runtime.NumCPU())
		for i := 0;  i < runtime.NumCPU(); i++ {
			go r.worker()
		}
	}

	r.Minimal = initialSet
	r.ddmin(initialSet, 2)

	return r, nil
}

func (r *Run) ddmin(set Set, n int) {
	subs, complements := r.split(set, n)

	// reduce to subset
	if nextSet := r.testSets(subs); nextSet != nil {
		r.ddmin(nextSet, 2)
		return
	}

	// reduce to complement
	if nextSet := r.testSets(complements); nextSet != nil {
		r.ddmin(nextSet, max(n-1, 2))
		return
	}

	// increase granularity
	if n < len(set) {
		r.ddmin(set, min(len(set), 2*n))
		return
	}

	// handle case where empty set of deltas causes failure of interest
	if empty := []int{}; r.Inp.Test(empty) == Failed {
		r.Minimal = empty
	}
}

func (r *Run) worker() {
	for set := range r.jobs {
		outcome := r.Inp.Test(set)
		r.results <- &Hist{Deltas: set, Out: outcome}
	}
}

func (r *Run) testSets(sets []Set) (failed Set) {
	if r.concurrent {
		// clear out old results
		for i := 0; i < len(r.results); i++ {
			<- r.results
		}

		for _, set := range sets {
			select {
			case r.jobs <- set:
				r.iterations++
			case hist := <- r.results:
				if hist.Out == Failed {
					r.Minimal = hist.Deltas
					return r.Minimal
				}
			}
		}

		for i := 0; i < len(r.results); i++ {
			if hist := <- r.results; hist.Out == Failed {
				r.Minimal = hist.Deltas
				return
			}
		}
		return nil
	}

	for _, set := range sets {
		result := r.Inp.Test(set)

		if r.keepHist {
			r.Hists = append(r.Hists, &Hist{Deltas: set, Out: result})
		} else {
			r.iterations++
		}

		if result == Failed {
			r.Minimal = set
			return set
		}
	}
	return nil
}

func (r *Run) IterCount() int {
	return max(len(r.Hists), r.iterations)
}

func (r *Run) split(set Set, n int) ([]Set, []Set) {
	size, remainder := len(set)/n, len(set)%n
	splits, complements := make([]Set, 0, n), make([]Set, 0, n)

	for i := 0; i < len(set)-remainder; i += size {
		complement := make(Set, 0, len(set)-size)
		complement = append(append(complement, set[:i]...), set[i+size:]...)
		split := set[i : i+size]
		if r.cacheTests {
			if hsh := split.hash(); !r.tested[hsh] {
				splits = append(splits, split)
				r.tested[hsh] = true
			}
			if hsh := complement.hash(); !r.tested[hsh] {
				complements = append(complements, complement)
				r.tested[hsh] = true
			}
			continue
		}
		splits = append(splits, split)
		complements = append(complements, complement)
	}

	if index := len(set) - remainder; index < len(set)-1 {
		splits = append(splits, set[index:])
		complements = append(complements, set[:index])
	}

	return splits, complements
}

func IntRange(n int) Set {
	r := make(Set, n)
	for i := 0; i < n; i++ {
		r[i] = i
	}
	return r
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
