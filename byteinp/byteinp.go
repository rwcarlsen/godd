
package byteinp

import (
  "bytes"
  "io"
  "io/ioutil"
  "sort"
  "github.com/rwcarlsen/godd"
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
  return &Word{words: bytes.Fields(data)}, nil
}

func (wi *Word) BuildInput(set godd.Set) []byte {
  sort.Ints([]int(set))

  inputWords := make([][]byte, len(set))
  for i, index := range set {
    inputWords[i] = wi.words[index]
  }

  input := bytes.Join(inputWords, []byte(" "))
  return append(input, byte('\n'))
}

func (wi *Word) Len() int {
  return len(wi.words)
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

func (wi *Char) BuildInput(set godd.Set) []byte {
  sort.Ints([]int(set))

  input := make([]byte, len(set))
  for i, index := range set {
    input[i] = wi.data[index]
  }

  return input
}

func (wi *Char) Len() int {
  return len(wi.data)
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

