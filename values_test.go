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

func TestValues_Equals(t *testing.T) {
	if !(Values{}).Equals(Values{}) {
		t.Errorf("Empty maps should be equal to each other")
	}
	if (Values{"a": 1}).Equals(Values{"a": true}) {
		t.Errorf("Maps with different values should not be equal")
	}
	if (Values{"a": 1}).Equals(Values{"a": 1, "b": 2}) {
		t.Errorf("Maps with a different number of keys should not be equal")
	}
}

func TestValues_Exclude(t *testing.T) {
	unsafe := Values{"a": 1, "b": 1}
	safe := unsafe.Exclude("a", "c")
	if len(safe) != 1 {
		t.Errorf("Unexpected length of safe Values: %d != 1", len(safe))
	}
	keys := safe.Keys()
	if len(keys) != 1 {
		t.Errorf("Unexpected length of safe keys: %d != 1", len(keys))
	}
	if keys[0] != "b" {
		t.Errorf("Unexpected safe key: %s != b", keys[0])
	}
}

func TestValues_Merge(t *testing.T) {
	a := Values{"a": 1, "b": 2}
	b := Values{"b": 3, "c": 4}
	c := a.Merge(b)

	if len(c) != 3 {
		t.Errorf("Unexpected length of c Values: %d != 3", len(c))
	}
	v, ok := c["c"].(int)
	if !ok {
		t.Fatal("Failed to convert the 'c' value to int")
	}
	if v != 4 {
		t.Errorf("Unexpected value of 'c': %d != 4", v)
	}

	// a should not be affected
	if _, exists := a["c"]; exists {
		t.Errorf("The original Values should not be modified")
	}
}

func TestValuesOf(t *testing.T) {
	// u := user{Email: "a@example.com", Name: "A"}
	// t.Error(ValuesOf(u))
}
