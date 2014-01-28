package godd

import (
	"testing"
)

type Tester Set

const TestSize = 10

func (bads Tester) Test(index Set) Outcome {
	found := true
	for _, bad := range bads {
		subfound := false
		for _, ti := range index {
			if ti == bad {
				subfound = true
				break
			}
		}
		found = found && subfound
	}
	if found {
		return Failed
	}
	return Passed
}

func (inp Tester) Len() int {
	return TestSize
}

var tests = []Set{
	Set{0},
	Set{TestSize - 1},
	Set{0, TestSize - 1},
	Set{1, 6, 8},
}

func TestMinFail(t *testing.T) {
	for i, set := range tests {
		tester := Tester(set)
		run, err := MinFail(tester)
		if err != nil {
			t.Errorf("set %v (%+v) failed: %v", i, set, err)
		}

		if len(run.MinFail) != len(set) {
			t.Errorf("set %v: expected %+v, got %+v", i, set, run.MinFail)
		}

		for i := range run.MinFail {
			if set[i] != run.MinFail[i] {
				t.Errorf("set %v: expected %+v, got %+v", i, set, run.MinFail)
			}
		}
		t.Logf("set %v (%v iter): %+v", i, len(run.Hists), run.MinFail)
		for _, h := range run.Hists {
			t.Log("    ", h)
		}
	}
}

func TestMinDiff(t *testing.T) {
	for i, set := range tests {
		tester := Tester(set)
		run, err := MinDiff(tester)
		if err != nil {
			t.Errorf("set %v (%+v) failed: %v", i, set, err)
		}

		if v := len(run.MinFail) - len(run.MaxPass); v != 1 {
			t.Errorf("set %v len(fail)-len(pass)=%v instead of 1", i, v)
			continue
		}

		t.Logf("set %v (%v iter) maxpass=%+v, minfail=%+v", i, len(run.Hists), run.MaxPass, run.MinFail)
		for _, h := range run.Hists {
			t.Log("    ", h)
		}
	}
}
