package types

import (
	"testing"
)

func TestNumeric(t *testing.T) {
	typ := Double()
	create, err := typ.Create(nil)
	if err != nil {
		t.Errorf("Unexpected error during DOUBLE Create(): %s", err)
	}
	if create != "DOUBLE PRECISION" {
		t.Errorf("Unexpected output of DOUBLE type: %s", create)
	}
}
