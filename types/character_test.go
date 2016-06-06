package types

import (
	"testing"
)

func TestCharacter(t *testing.T) {
	columnType := Varchar(32)
	create, err := columnType.Create(nil) // No dialect needed
	if err != nil {
		t.Errorf("Unexpected error during VARCHAR Create(): %s", err)
	}
	if create != "VARCHAR(32)" {
		t.Errorf("Unexpected output of VARCHAR type: %s", create)
	}
}
