
package godd

import (
	"testing"
  "fmt"
)

type TestInput int

func (inp TestInput) Test(index Set) bool {
	tot := len(index) + 1
	return tot % 2 == 0
}

func (inp TestInput) Len() int {
	return int(inp)
}

type TestInput2 []int

func (inp TestInput2) Test(index Set) bool {
  for _, failPart := range inp {
    found := false
    for _, v := range index {
      if v == failPart {
        found = true
        break
      }
    }
    if !found {
      return true
    }
  }
  return false
}

func (inp TestInput2) Len() int {
	return 50
}

func TestMinFail(t *testing.T) {
	test1(t, TestInput(12))
	test1(t, TestInput(8))
	test1(t, TestInput(2))
	test2(t, TestInput2([]int{1, 5, 6, 19, 47}))
}

func test1(t *testing.T, inp TestInput) {
	t.Logf("\n--------- input length %v ----------\n", inp.Len())

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
		t.Logf("hist %v (%v): %v \n", i, result, hist.DeltaInd)
	}

	t.Logf("minimal failing input (%v iterations): %v\n", len(run.Hists), run.Minimal)

	if len(run.Minimal) != 2 {
		t.Errorf("FAILED:: minimal output: got %v, expected [# #]", run.Minimal)
	}
}

func test2(t *testing.T, inp TestInput2) {
	t.Logf("\n--------- input length %v ----------\n", inp.Len())

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
		t.Logf("hist %v (%v): %v \n", i, result, hist.DeltaInd)
	}

	t.Logf("minimal failing input (%v iterations): %v\n", len(run.Hists), run.Minimal)

	if fmt.Sprint(run.Minimal) != fmt.Sprint(inp) {
		t.Errorf("FAILED:: minimal output: got %v, expected %v", run.Minimal, inp)
	}
}
