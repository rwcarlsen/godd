package godd

import (
	"testing"
)

type Tester Set

const TestSize = 1000

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
	Set{2, 356, 358},
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
		t.Logf("set %v: expected %+v, got %+v", i, set, run.MinFail)
	}
}

func TestMinDiff(t *testing.T) {
	for i, set := range tests {
		tester := Tester(set)
		run, err := MinDiff(tester)
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
		t.Logf("set %v minfail: expected %+v, got %+v", i, set, run.MinFail)
		t.Logf("set %v minpass: %+v", i, run.MinPass)
	}
}
