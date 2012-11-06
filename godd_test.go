
package godd

import (
	"testing"
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
	return str + "1"
}

func TestMinFail(t *testing.T) {
	testN(t, 12)
	testN(t, 8)
	testN(t, 3)
	testN(t, 2)
}

func testN(t *testing.T, n int) {
	inp := &TestInput{n: n}
	run := &Run{Inp: inp}

	run.MinFail()

	t.Logf("\n--------- input length %v ----------\n", n)
	for i, hist := range run.Hists {
		result := "PASS"
		if !hist.Passed {
			result = "FAIL"
		}
		t.Logf("hist %v (%v): %v %v\n", i, result, hist.Set, inp.Compose(hist.Set))
	}

	t.Logf("minimal failing input: %v\n", inp.Compose(run.Minimal))

	if len(run.Minimal) != 2 {
		t.Errorf("minimal output: got %v, expected [# #]", run.Minimal)
	}
}
