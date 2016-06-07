package types

import (
	"testing"
)

func TestDatetime(t *testing.T) {
	datatype := Datetime()
	create, err := datatype.Create(nil) // No dialect needed
	if err != nil {
		t.Errorf("Unexpected error during DATETIME Create(): %s", err)
	}
	if create != "DATETIME" {
		t.Errorf("Unexpected output of DATETIME type: %s", create)
	}
}
