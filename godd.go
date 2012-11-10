package godd

import (
	"errors"
	"os"
)

var buf = buffer(false)

type buffer bool

func (b buffer) Write(p []byte) (n int, err error) {
	if b {
		return os.Stdout.Write(p)
	}
	return len(p), nil
}

type Set []int

func (s Set) hash() string {
	h := make([]byte, len(s))
	for i, v := range s {
		h[i] = byte(v)
	}
	return string(h)
}

type Input interface {
	Passes(Set) bool
	Len() int
}

type Hist struct {
	DeltaInd Set
	Passed   bool
}

type Run struct {
	Inp     Input
	Minimal Set
	Hists   []*Hist
	tested  map[string]bool
}

func MinFail(inp Input) (*Run, error) {
	r := &Run{Inp: inp}
	r.tested = make(map[string]bool)
	initialSet := IntRange(inp.Len())

	if passed := inp.Passes(initialSet); passed {
		return nil, errors.New("godd: Test passes with all deltas applied")
	}

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
	if empty := []int{}; r.Inp.Passes(empty) == false {
	  r.Minimal = empty
	}
}

func (r *Run) testSets(sets []Set) (failed Set) {
	for _, set := range sets {
		if r.tested[set.hash()] {
			continue
		}
		r.tested[set.hash()] = true

		passed := r.Inp.Passes(set)
		r.Hists = append(r.Hists, &Hist{DeltaInd: set, Passed: passed})

		if !passed {
			r.Minimal = set
			return set
		}
	}
	return nil
}

func split(set Set, n int) ([]Set, []Set) {
	size, remainder := len(set)/n, len(set)%n

	splits := make([]Set, n)
	complements := make([]Set, n)
	count := 0
	for i := 0; i < len(set)-remainder; i += size {
		splits[count] = set[i:i+size]
		complement := append(Set{}, set[:i]...)
		complement = append(complement, set[i+size:]...)
		complements[count] = complement
		count++
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
