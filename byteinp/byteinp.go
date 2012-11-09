package byteinp

import (
	"bytes"
	"github.com/rwcarlsen/godd"
	"io"
	"io/ioutil"
	"sort"
)

type Builder interface {
	BuildInput(godd.Set) []byte
	Len() int
}

type Tester interface {
	Test(input []byte) bool
}

type Word struct {
	words [][]byte
}

func ByWord(r io.Reader) (*Word, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	words := splitAny(data, "\n\t ")
	return &Word{words: words}, nil
}

func (wb *Word) BuildInput(set godd.Set) []byte {
	sort.Ints([]int(set))

	inputWords := make([][]byte, len(set))
	for i, index := range set {
		inputWords[i] = wb.words[index]
	}

	input := bytes.Join(inputWords, []byte(""))
	return input
}

func (wb *Word) Len() int {
	return len(wb.words)
}

type Char struct {
	data []byte
}

func ByChar(r io.Reader) (*Char, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Char{data: data}, nil
}

func (cb *Char) BuildInput(set godd.Set) []byte {
	sort.Ints([]int(set))

	input := make([]byte, len(set))
	for i, index := range set {
		input[i] = cb.data[index]
	}

	return input
}

func (cb *Char) Len() int {
	return len(cb.data)
}

type Line struct {
	lines [][]byte
}

func ByLine(r io.Reader) (*Line, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(data, []byte("\n"))
	return &Line{lines: lines}, nil
}

func (lb *Line) BuildInput(set godd.Set) []byte {
	sort.Ints([]int(set))

	inputLines := make([][]byte, len(set))
	for i, index := range set {
		inputLines[i] = lb.lines[index]
	}

	input := bytes.Join(inputLines, []byte("\n"))
	return append(input, byte('\n'))
}

func (lb *Line) Len() int {
	return len(lb.lines)
}

type TestCase struct {
	T Tester
	B Builder
}

func (t *TestCase) Passes(set godd.Set) bool {
	input := t.B.BuildInput(set)
	return t.T.Test(input)
}

func (t *TestCase) Len() int {
	return t.B.Len()
}

// splitAny splits a byte slice at every occurence of any char in chars -
// collapsing consecutive char delimiters onto the preceding split chunk.
func splitAny(s []byte, chars string) [][]byte {
	beg := 0
	chunks := [][]byte{}
	for i, b := range s {
		for _, c := range chars {
			if b == byte(c) {
				end := i + 1
				nextNotWhite := bytes.IndexAny(s[end:end+1], chars) == -1
				if nextNotWhite {
					chunks = append(chunks, s[beg:end])
					beg = end
				}
				break
			}
		}
	}
	return append(chunks, s[beg:])
}
