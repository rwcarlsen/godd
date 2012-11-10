package godd

import (
	"math/rand"
	"testing"
	"sort"
)

type TestInput int

func (inp TestInput) Passes(index Set) bool {
	tot := len(index) + 1
	return tot%2 == 0
}

func (inp TestInput) Len() int {
	return int(inp)
}

type TestInput2 []int

func (inp TestInput2) Passes(index Set) bool {
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

func (_ TestInput2) Len() int {
	return 20000
}

func TestMinFail(t *testing.T) {
	test1(t, TestInput(12))
	test1(t, TestInput(8))
	test1(t, TestInput(2))

	inp := make([]int, 200)
	for i := 0; i < 200; i++ {
		inp[i] = rand.Intn(20000)
	}
	sort.Ints(inp)
	test2(t, TestInput2(inp))
}

func test1(t *testing.T, inp TestInput) {
	t.Logf("\n--------- input length %v ----------\n", inp.Len())

	run, err := MinFail(inp)

	if err != nil {
		t.Errorf("FAILED: %v", err)
		return
	}

	t.Logf("minimal failing input (%v iterations): %v\n", len(run.Hists), run.Minimal)

	if len(run.Minimal) != 0 {
		t.Errorf("FAILED:: minimal output: got %v, expected []", run.Minimal)
	}
}

func test2(t *testing.T, inp TestInput2) {
	t.Logf("\n--------- input length %v ----------\n", inp.Len())

	run, err := MinFail(inp)

	if err != nil {
		t.Errorf("FAILED: %v", err)
		return
	}

	t.Logf("minimal failing input (%v iterations): %v\n", len(run.Hists), run.Minimal)

	//if fmt.Sprint(run.Minimal) != fmt.Sprint(inp) {
	//	t.Errorf("FAILED:: minimal output: expected %v", inp)
	//}
}
