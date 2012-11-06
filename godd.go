
package godd

import (
  "fmt"
  "log"
)

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
  log.Println("--------------------- begin test ----------------------")
  r.tested = make(map[string]bool)
	r.ddmin(intRange(r.Inp.Len()), 2)
	return nil
}

func (r *Run) ddmin(set Set, n int) {
	subs, complements := split(set, n)
  log.Println("--------- recurse ------------")
  log.Println("subs: ", subs)
  log.Println("complements: ", complements)

	// reduce to subset
  if nextSet := r.testSets(subs); nextSet != nil {
      log.Println("reducing to subset...")
			r.ddmin(nextSet, 2)
      return
	}

	// reduce to complement
  if nextSet := r.testSets(complements); nextSet != nil {
      log.Println("reducing to complement...")
			r.ddmin(nextSet, max(n-1, 2))
      return
	}

	// increase granularity
	if n < len(set) {
    log.Println("increase granularity...")
		r.ddmin(set, min(len(set), 2 * n))
	}
}

func (r *Run) testSets(sets []Set) (failed Set) {
	for _, set := range sets {
    if r.tested[set.hashable()] {
      continue
    }

    passed := r.Inp.Test(set)

    r.tested[set.hashable()] = true
		r.Hists = append(r.Hists, &Hist{DeltaInd: set, Passed: passed})

		if !passed {
			r.Minimal = set
      return set
		}
	}
  return nil
}

func split(set Set, n int) ([]Set, []Set) {
	size, remainder := len(set) / n, len(set) % n

	splits := []Set{}
	complements := []Set{}
	for i := 0; i < len(set) - remainder; i += size {
		splits = append(splits, set[i:i+size])
		complement := append(Set{}, set[:i]...)
		complement = append(complement, set[i+size:]...)
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

