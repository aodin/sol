package types

import (
	"testing"
)

func TestBoolean(t *testing.T) {
	datatype := Boolean()
	create, err := datatype.Create(nil) // No dialect needed
	if err != nil {
		t.Errorf("Unexpected error during BOOLEAN Create(): %s", err)
	}
	if create != "BOOLEAN" {
		t.Errorf("Unexpected output of BOOLEAN type: %s", create)
	}

	datatype = datatype.NotNull().Unique().Default(false)
	create, err = datatype.Create(nil) // No dialect needed
	if err != nil {
		t.Errorf(
			"Unexpected error during BOOLEAN Create() with options: %s", err,
		)
	}
	if create != "BOOLEAN NOT NULL UNIQUE DEFAULT false" {
		t.Errorf("Unexpected output of BOOLEAN type with options: %s", create)
	}
}
