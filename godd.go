
package godd

import (
	"fmt"
)

type Input interface {
	Test([]int) bool
	Len() int
}

func MinFail(inp Input) []int {
	m := &minner{inp: inp}
	return m.ddmin(intRange(inp.Len()), 2)
}

type minner struct {
	inp Input
}

func (m *minner) ddmin(set []int, n int) []int {
	fmt.Println("recurse --------------")
	// reduce to subset
	subs, complements := split(set, n)
	for _, sub := range subs {
		if ! m.inp.Test(sub) {
			m.hist = append(m.hist, sub)
			return m.ddmin(sub, 2)
		}
	}

	// reduce to complement
	for _, comp := range complements {
		if ! m.inp.Test(comp) {
			m.hist = append(m.hist, comp)
			return m.ddmin(comp, max(n-1, 2))
		}
	}

	// increase granularity
	if n < len(set) {
		return m.ddmin(set, min(len(set), 2 * n))
	}

	return set
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

func split(set []int, n int) ([][]int, [][]int) {
	size, remainder := len(set) / n, len(set) % n

	splits := [][]int{}
	complements := [][]int{}
	for i := 0; i < n - remainder; i += size {
		splits = append(splits, set[i:i+size])
		complement := append([]int{}, set[:i]...)
		complement = append(complement, set[i+1:]...)
		complements = append(complements, complement)
	}

	index := len(set)-1-remainder
	splits = append(splits, set[index:])
	complements = append(complements, set[:index])

	return splits, complements
}

func intRange(n int) []int {
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = i
	}
	return r
}

