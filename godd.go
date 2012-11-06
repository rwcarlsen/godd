
package godd

import "fmt"

type Set []int

func (s Set) hashable() string {
  return fmt.Sprint(s)
}

type Input interface {
	Test(Set) bool
	Len() int
}

type Hist struct {
	DeltaInd Set
	Passed bool
}

type Run struct {
	Inp Input
	Minimal Set
	Hists []*Hist
  tested map[string]bool
}

func (r *Run) MinFail() error {
  r.tested = make(map[string]bool)
	r.ddmin(intRange(r.Inp.Len()), 2)
	return nil
}

func testSets(sets []Set, n int) (passed bool) {
	for _, set := range sets {
    if r.tested[set.hashable()] {
      continue
    }
    passed = r.Inp.Test(sub)
    r.tested[sub.hashable()] = true
		r.Hists = append(r.Hists, &Hist{DeltaInd: sub, Passed: passed})
		if !passed {
			r.Minimal = sub
      return false
		}
	}
  return true
}

func (r *Run) ddmin(set Set, n int) {
	subs, complements := split(set, n)

	// reduce to subset
	for _, sub := range subs {
    if r.tested[sub.hashable()] {
      continue
    }
    passed := r.Inp.Test(sub)
    r.tested[sub.hashable()] = true
		r.Hists = append(r.Hists, &Hist{DeltaInd: sub, Passed: passed})
		if !passed {
			r.Minimal = sub
			r.ddmin(sub, 2)
			return
		}
	}

	// reduce to complement
	for _, comp := range complements {
    if r.tested[comp.hashable()] {
      continue
    }
		passed := r.Inp.Test(comp)
    r.tested[comp.hashable()] = true
		r.Hists = append(r.Hists, &Hist{DeltaInd: comp, Passed: passed})
		if ! passed {
			r.Minimal = comp
			r.ddmin(comp, max(n-1, 2))
			return
		}
	}

	// increase granularity
	if n < len(set) {
		r.ddmin(set, min(len(set), 2 * n))
		return
	}
}

func split(set Set, n int) ([]Set, []Set) {
	size, remainder := len(set) / n, len(set) % n

	splits := []Set{}
	complements := []Set{}
	for i := 0; i < len(set) - remainder; i += size {
		splits = append(splits, set[i:i+size])
		complement := append(Set{}, set[:i]...)
		complement = append(complement, set[i+1:]...)
		complements = append(complements, complement)
	}

	if index := len(set)-remainder; index < len(set)-1 {
		splits = append(splits, set[index:])
		complements = append(complements, set[:index])
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

