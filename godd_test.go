package godd

import (
	"math/rand"
	"sort"
	"testing"
)

type TestInput []int

var staticInp []int

func init() {
	staticInp = make([]int, 150)
	for i := 0; i < 150; i++ {
		staticInp[i] = rand.Intn(10000)
	}
	sort.Ints(staticInp)
}

func (inp TestInput) Test(index Set) Outcome {
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

func (_ TestInput) Len() int {
	return 10000
}

func TestMinFail(t *testing.T) {
	inp := TestInput(staticInp)
	run, err := MinFail(inp, CcacheTests)
	if err != nil {
		t.Errorf("FAILED: %v", err)
		return
	}
	t.Logf("minimal failing input (%v iterations, len %v): %v\n",
			run.IterCount(), len(run.Minimal), run.Minimal)
}

func BenchmarkMinFail_NoCacheNoHist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inp := TestInput(staticInp)
		run, _ := MinFail(inp, 0)
		b.Logf("%v iterations", run.IterCount())
	}
}

func BenchmarkMinFail_NoCacheHist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inp := TestInput(staticInp)
		run, _ := MinFail(inp, CkeepHist)
		b.Logf("%v iterations", run.IterCount())
	}
}

func BenchmarkMinFail_CacheNoHist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inp := TestInput(staticInp)
		run, _ := MinFail(inp, CcacheTests)
		b.Logf("%v iterations", run.IterCount())
	}
}

func BenchmarkMinFail_CacheHist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inp := TestInput(staticInp)
		run, _ := MinFail(inp, CcacheTests|CkeepHist)
		b.Logf("%v iterations", run.IterCount())
	}
}

func BenchmarkMinFail_Concurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inp := TestInput(staticInp)
		run, _ := MinFail(inp, CcacheTests|Cconcurrent)
		b.Logf("%v iterations", run.IterCount())
	}
}
