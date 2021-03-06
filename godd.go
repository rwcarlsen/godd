package godd

import (
	"errors"
	"fmt"
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
	panic("using invalid outcome")
}

const (
	Passed Outcome = iota
	Failed
	Undetermined
)

type Set []int

func Sub(a, b Set) Set {
	m := make(map[int]struct{}, len(b))
	for _, v := range b {
		m[v] = struct{}{}
	}
	c := make(Set, 0, len(a))
	for _, v := range a {
		if _, ok := m[v]; !ok {
			c = append(c, v)
		}
	}
	return c
}

type Input interface {
	Test(Set) Outcome
	Len() int
}

type Hist struct {
	DeltaInd Set
	Out      Outcome
}

func (h *Hist) String() string {
	return fmt.Sprintf("%v: %v", h.DeltaInd, h.Out)
}

type Run struct {
	Inp     Input
	MinFail Set
	MaxPass Set
	Hists   []*Hist
}

func MinFail(inp Input) (*Run, error) {
	r := &Run{Inp: inp}
	initialSet := intRange(inp.Len())

	if inp.Test(initialSet) != Failed {
		return nil, errors.New("godd: Test passes with all deltas applied")
	}

	r.ddmin(initialSet, 2)
	return r, nil
}

func MinDiff(inp Input) (*Run, error) {
	r := &Run{Inp: inp}
	initialSet := intRange(inp.Len())

	if inp.Test(initialSet) != Failed {
		return nil, errors.New("godd: Test passes with all deltas applied")
	} else if inp.Test(Set{}) != Passed {
		return nil, errors.New("godd: Test does not pass with no deltas applied")
	}

	r.dd(Set{}, initialSet, 2)
	return r, nil
}

func (r *Run) dd(passing, failing Set, n int) {
	r.MinFail = failing
	r.MaxPass = passing
	set := Sub(failing, passing)
	if n > len(set) {
		return
	}

	subs, complements := split(passing, failing, n)

	// reduce to subset
	if nextSet := r.testSets(subs, Failed); nextSet != nil {
		r.dd(passing, nextSet, 2)
		return
	}

	// increase to complement
	if nextSet := r.testSets(complements, Passed); nextSet != nil {
		r.dd(nextSet, failing, 2)
		return
	}

	// increase to subset
	if nextSet := r.testSets(subs, Passed); nextSet != nil {
		r.dd(nextSet, failing, max(n-1, 2))
		return
	}

	// reduce to complement
	if nextSet := r.testSets(complements, Failed); nextSet != nil {
		r.dd(passing, nextSet, max(n-1, 2))
		return
	}

	// increase granularity
	if n < len(set) {
		r.dd(passing, failing, min(len(set), 2*n))
		return
	}

	// otherwise, done
}

func (r *Run) ddmin(set Set, n int) {
	r.MinFail = set
	if n > len(set) {
		return
	}

	subs, complements := splitMin(set, n)

	// reduce to subset
	if nextSet := r.testSets(subs, Failed); nextSet != nil {
		r.ddmin(nextSet, 2)
		return
	}

	// reduce to complement
	if nextSet := r.testSets(complements, Failed); nextSet != nil {
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
		r.MinFail = empty
	}

	// otherwise, done
}

func (r *Run) testSets(sets []Set, expected Outcome) (matched Set) {
	for _, set := range sets {
		result := r.Inp.Test(set)
		r.Hists = append(r.Hists, &Hist{DeltaInd: set, Out: result})

		if expected == result {
			return set
		}
	}
	return nil
}

func splitMin(set Set, n int) ([]Set, []Set) {
	size, remainder := len(set)/n, len(set)%n
	splits, complements := make([]Set, n), make([]Set, n)

	count := 0
	for i := 0; i < len(set)-remainder; i += size {
		splits[count] = append(splits[count], set[i:i+size]...)
		complement := make(Set, 0, len(set)-size)
		complement = append(append(complement, set[:i]...), set[i+size:]...)
		complements[count] = complement
		count++
	}

	index := len(set) - remainder
	splits[n-1] = append(splits[n-1], set[index:]...)

	if len(complements[n-1]) > 0 {
		complements[n-1] = complements[n-1][:len(complements[n-1])-remainder]
	}

	return splits, complements
}

func split(passing, failing Set, n int) ([]Set, []Set) {
	set := Sub(failing, passing)
	size, remainder := len(set)/n, len(set)%n
	splits, complements := make([]Set, n), make([]Set, n)

	count := 0
	for i := 0; i < len(set)-remainder; i += size {
		splits[count] = append(splits[count], passing...)
		splits[count] = append(splits[count], set[i:i+size]...)
		complements[count] = Sub(failing, set[i:i+size])
		count++
	}

	index := len(set) - remainder
	splits[n-1] = append(splits[n-1], set[index:]...)

	if len(complements[n-1]) > 0 {
		complements[n-1] = Sub(complements[n-1], set[index:])
	}

	return splits, complements
}

func intRange(n int) Set {
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
