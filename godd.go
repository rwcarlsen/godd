
package godd

import "fmt"

type Input interface {
	Test([]int) bool
	Len() int
}

type Hist struct {
	Set []int
	Passed bool
}

type Run struct {
	Inp Input
	Minimal []int
	Hists []*Hist
}

func (r *Run) MinFail() error {
	r.ddmin(intRange(r.Inp.Len()), 2)
	return nil
}


func (r *Run) ddmin(set []int, n int) {
	fmt.Println("recurse ----------")
	subs, complements := split(set, n)
	fmt.Println("subs: ", subs)
	fmt.Println("complements: ", complements)

	// reduce to subset
	for _, sub := range subs {
		passed := r.Inp.Test(sub)
		r.Hists = append(r.Hists, &Hist{Set: sub, Passed: passed})
		hist := r.Hists[len(r.Hists)-1]
		fmt.Printf("hist (%v): %v \n", passed, hist.Set)
		if !passed {
			r.Minimal = sub
			r.ddmin(sub, 2)
			return
		}
	}

	// reduce to complement
	for _, comp := range complements {
		passed := r.Inp.Test(comp)
		r.Hists = append(r.Hists, &Hist{Set: comp, Passed: passed})
		hist := r.Hists[len(r.Hists)-1]
		fmt.Printf("hist (%v): %v \n", passed, hist.Set)
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

func split(set []int, n int) ([][]int, [][]int) {
	size, remainder := len(set) / n, len(set) % n

	splits := [][]int{}
	complements := [][]int{}
	for i := 0; i < len(set) - remainder; i += size {
		splits = append(splits, set[i:i+size])
		complement := append([]int{}, set[:i]...)
		complement = append(complement, set[i+1:]...)
		complements = append(complements, complement)
	}

	if index := len(set)-remainder; index < len(set)-1 {
		splits = append(splits, set[index:])
		complements = append(complements, set[:index])
	}

	return splits, complements
}

func intRange(n int) []int {
	r := make([]int, n)
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

