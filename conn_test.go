package sol

import (
	"database/sql"
	"testing"
)

var _ executer = &sql.DB{}
var _ executer = &sql.Tx{}

func TestConn(t *testing.T) {
	// Creating a Must() connection should not modify the original connection
	c := &conn{}
	d := c.Must()
	if !d.panicky {
		t.Errorf("Expected Must() connection to have panicky = true")
	}
	if c.panicky {
		t.Errorf("Original connection should not be modified by Must()")
	}
}
