
package byteinp

import (
  "testing"
  "bytes"
  )

func TestSplitAny(t *testing.T) {
  data := []byte("this  is\t\t a line\t that has\t\t lots of  white space")
  splits := splitAny(data, "\n\t ")
  joined := bytes.Join(splits, []byte(","))
  t.Log(string(joined))
}
