package sol

import (
	"encoding/json"
	"testing"
)

func TestValues(t *testing.T) {
	values := Values{"c": []byte("bytes")}
	b, err := json.Marshal(values)
	if err != nil {
		t.Fatalf("Unexpected error while marshaling: %s", err)
	}
	if string(b) != `{"c":"bytes"}` {
		t.Errorf(`Unexpected JSON marshal: %s != {"c":"bytes"}`, b)
	}
}

func TestValues_Exclude(t *testing.T) {
	unsafe := Values{"a": 1, "b": 1}
	safe := unsafe.Exclude("b", "c")
	if len(safe) != 1 {
		t.Errorf("Unexpected length of safe Values: %d != 1", len(safe))
	}
	keys := safe.Keys()
	if len(keys) != 1 {
		t.Errorf("Unexpected length of safe keys: %d != 1", len(keys))
	}
	if keys[0] != "a" {
		t.Errorf("Unexpected safe key: %s != a", keys[0])
	}
}

func TestValuesOf(t *testing.T) {
	// u := user{Email: "a@example.com", Name: "A"}
	// t.Error(ValuesOf(u))
}
