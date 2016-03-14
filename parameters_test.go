package sol

import "testing"

func TestParameters_Add(t *testing.T) {
	ps := Params()
	ps.Add(1)

	if len(*ps) != 1 {
		t.Fatalf("Unexpected length of parameters: %d != 1", len(*ps))
	}
}
