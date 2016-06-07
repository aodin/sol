package types

import (
	"testing"
)

func TestNumeric(t *testing.T) {
	datatype := Double()
	create, err := datatype.Create(nil) // No dialect needed
	if err != nil {
		t.Errorf("Unexpected error during DOUBLE Create(): %s", err)
	}
	if create != "DOUBLE PRECISION" {
		t.Errorf("Unexpected output of DOUBLE type: %s", create)
	}
}
