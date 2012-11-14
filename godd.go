package godd

import (
	"errors"
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
	iterations int
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

	r.Minimal = initialSet
	r.ddmin(initialSet, 2)
	return r, nil
}

func (r *Run) ddmin(set Set, n int) {
	subs, complements := split(set, n)

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

func (r *Run) workerListen(reset chan bool) {

}

func (r *Run) testSets(sets []Set) (failed Set) {
	for _, set := range sets {
		if r.cacheTests {
			if r.tested[set.hash()] {
				continue
			}
			r.tested[set.hash()] = true
		}

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

func split(set Set, n int) ([]Set, []Set) {
	size, remainder := len(set)/n, len(set)%n
	splits, complements := make([]Set, n), make([]Set, n)

	count := 0
	for i := 0; i < len(set)-remainder; i += size {
		splits[count] = set[i : i+size]
		complement := make(Set, 0, len(set)-size)
		complement = append(append(complement, set[:i]...), set[i+size:]...)
		complements[count] = complement
		count++
	}

	if index := len(set) - remainder; index < len(set)-1 {
		splits[n-1] = set[index:]
		complements[n-1] = set[:index]
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
