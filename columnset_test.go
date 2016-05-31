package sol

import "testing"

func TestColumns(t *testing.T) {
	columns := Columns()
	example := ColumnElem{name: "example"}

	added, err := columns.Add(example)
	if err != nil {
		t.Fatalf("sol: adding a column to a set should not error: %s", err)
	}
	if len(added.order) != 1 {
		t.Errorf(
			"sol: unexpected length of ColumnSet: 1 != %d",
			len(added.order),
		)
	}
	if !added.Has("example") {
		t.Errorf("sol: ColumnSet should have a column named 'example")
	}
	if added.Has("test") {
		t.Errorf("sol: ColumnSet should not have a column named 'test")
	}
}

func TestUniqueColumns(t *testing.T) {
	columns := UniqueColumns()
	example := ColumnElem{name: "example"}

	added, err := columns.Add(example)
	if err != nil {
		t.Fatalf("sol: adding a column to a set should not error: %s", err)
	}
	_, err = added.Add(example)
	if err == nil {
		t.Fatalf("sol: adding a duplicate column to a unqie set should error")
	}

	rejected := columns.Reject("example")
	if len(rejected.order) != 0 {
		t.Errorf("sol: Reject should have removed the only column")
	}
}
