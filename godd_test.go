
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

	fmt.Println(inp.Compose(index), ": pass=", tot%2==0)
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
	inp := &TestInput{n: 8}

	minimal := MinFail(inp)
	fmt.Println(inp.Compose(minimal))
}
