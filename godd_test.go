package godd

import (
	"math/rand"
	"sort"
	"testing"
)

type TestInput2 []int

var staticInp []int

func init() {
	staticInp = make([]int, 200)
	for i := 0; i < 200; i++ {
		staticInp[i] = rand.Intn(15000)
	}
	sort.Ints(staticInp)
}

func (inp TestInput2) Test(index Set) Outcome {
	for _, failPart := range inp {
		found := false
		for _, v := range index {
			if v == failPart {
				found = true
				break
			}
		}
		if !found {
			return Passed
		}
	}
	return Failed
}

func (_ TestInput2) Len() int {
	return 15000
}

func TestMinFail(t *testing.T) {
	inp := TestInput2(staticInp)
	run, err := MinFail(inp, Cdefault)
	if err != nil {
		t.Errorf("FAILED: %v", err)
		return
	}
	t.Logf("minimal failing input (%v iterations, len %v): %v\n", len(run.Hists), len(run.Minimal), run.Minimal)
}

func test(t *testing.T, inp TestInput2) {
}

func BenchmarkMinFail_NoCacheNoHist(b *testing.B) {
  for i := 0; i < b.N; i++ {
    inp := TestInput2(staticInp)
    MinFail(inp, 0)
  }
}

func BenchmarkMinFail_NoCacheHist(b *testing.B) {
  for i := 0; i < b.N; i++ {
    inp := TestInput2(staticInp)
    MinFail(inp, CkeepHist)
  }
}

func BenchmarkMinFail_CacheNoHist(b *testing.B) {
  for i := 0; i < b.N; i++ {
    inp := TestInput2(staticInp)
    MinFail(inp, CcacheTests)
  }
}

func BenchmarkMinFail_CacheHist(b *testing.B) {
  for i := 0; i < b.N; i++ {
    inp := TestInput2(staticInp)
    MinFail(inp, CcacheTests | CkeepHist)
  }
}
