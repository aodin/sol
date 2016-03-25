package sol

import "testing"

func TestConn(t *testing.T) {
	// Creating a Must() connection should not modify the original connection
	c := &DBConn{}
	d := c.Must()
	if !d.panicky {
		t.Errorf("Expected Must() connection to have panicky = true")
	}
	if c.panicky {
		t.Errorf("Original connection should not be modified by Must()")
	}
}
