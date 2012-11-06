
package godd

import (
	"testing"
	"fmt"
)

type TestInput struct {
	n int
}

func (inp *TestInput) Test(index []int) bool {
	tot := len(index) + 1
	return tot % 2 == 0
}

func (inp *TestInput) Len() int {
	return inp.n
}

func (inp *TestInput) Compose(index []int) string {
	str := ""
	for _ = range index {
		str += "1+"
	}
	return str + "   1"
}

func TestMinFail(t *testing.T) {
	inp := &TestInput{n: 8}
	run := &Run{Inp: inp}

	run.MinFail()

	for i, hist := range run.Hists {
		result := "PASS"
		if !hist.Passed {
			result = "FAIL"
		}
		fmt.Printf("hist %v (%v): %v %v\n", i, result, hist.Set, inp.Compose(hist.Set))
	}

	fmt.Printf("minimal failing input: %v\n", inp.Compose(run.Minimal))
}
