
package godd

import (
	"testing"
)

type TestInput struct {
	n int
}

func (inp *TestInput) Test(index Set) bool {
	tot := len(index) + 1
	return tot % 2 == 0
}

func (inp *TestInput) Len() int {
	return inp.n
}

func (inp *TestInput) Compose(index Set) string {
	str := ""
	for _ = range index {
		str += "1+"
	}
	return str + "1"
}

func TestMinFail(t *testing.T) {
	testN(t, 12)
	testN(t, 8)
	testN(t, 2)
}

func testN(t *testing.T, n int) {
	t.Logf("\n--------- input length %v ----------\n", n)

	inp := &TestInput{n: n}
	run, err := MinFail(inp)

	if err != nil {
    t.Errorf("FAILED: %v", err)
    return
  }

	for i, hist := range run.Hists {
		result := "PASS"
		if !hist.Passed {
			result = "FAIL"
		}
		t.Logf("hist %v (%v): %v %v\n", i, result, hist.DeltaInd, inp.Compose(hist.DeltaInd))
	}

	t.Logf("minimal failing input: %v\n", inp.Compose(run.Minimal))

	if len(run.Minimal) != 2 {
		t.Errorf("FAILED:: minimal output: got %v, expected [# #]", run.Minimal)
	}
}
