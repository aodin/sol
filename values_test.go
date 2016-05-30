package sol

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
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
	var out Values
	var err error

	// Test map types
	values := Values{"a": 1}
	if out, err = ValuesOf(values); err != nil {
		t.Errorf("Unexpected error from ValuesOf() with Values: %s", err)
	}
	if !reflect.DeepEqual(values, out) {
		t.Errorf("Unexpected values from ValuesOf() with Values: %+v", out)
	}
	if out, err = ValuesOf(&values); err != nil {
		t.Errorf("Unexpected error from ValuesOf() with *Values: %s", err)
	}
	if !reflect.DeepEqual(values, out) {
		t.Errorf("Unexpected values from ValuesOf() with *Values: %+v", out)
	}

	attrs := map[string]interface{}{"b": 2}
	if out, err = ValuesOf(attrs); err != nil {
		t.Errorf("Unexpected error from ValuesOf() with map: %s", err)
	}
	if !reflect.DeepEqual(Values(attrs), out) {
		t.Errorf("Unexpected values from ValuesOf() with map: %+v", out)
	}
	if out, err = ValuesOf(&attrs); err != nil {
		t.Errorf("Unexpected error from ValuesOf() with *map: %s", err)
	}
	if !reflect.DeepEqual(Values(attrs), out) {
		t.Errorf("Unexpected values from ValuesOf() with *map: %+v", out)
	}

	// The following types are declared in fields_test
	embed := embedded{
		Serial: Serial{ID: uint64(20)},
		Name:   "Object",
	}

	if out, err = ValuesOf(embed); err != nil {
		t.Errorf("Unexpected error from ValuesOf() with struct: %s", err)
	}

	expected := Values{
		"id":        uint64(20),
		"Name":      "Object",
		"UpdatedAt": (*time.Time)(nil),
	}
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Unexpected values from ValuesOf() with *struct: %+v", out)
	}
	if out, err = ValuesOf(embed); err != nil {
		t.Errorf("Unexpected error from ValuesOf() with *struct: %s", err)
	}
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Unexpected values from ValuesOf() with *struct: %+v", out)
	}

	now := time.Now()
	embed.Serial.ID = uint64(0)
	embed.Timestamp.CreatedAt = now

	if out, err = ValuesOf(embed); err != nil {
		t.Errorf("Unexpected error from ValuesOf() with struct: %s", err)
	}

	expected = Values{
		"Name":      "Object",
		"CreatedAt": now,
		"UpdatedAt": (*time.Time)(nil),
	}
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Unexpected values from ValuesOf() with *struct: %+v", out)
	}
}

var emptyValues = []bool{
	isEmptyValue(reflect.ValueOf(0)),
	isEmptyValue(reflect.ValueOf("")),
	isEmptyValue(reflect.ValueOf(false)),
	isEmptyValue(reflect.ValueOf(0.0)),
	isEmptyValue(reflect.ValueOf(time.Time{})),
}

var nonEmptyValues = []bool{
	isEmptyValue(reflect.ValueOf(1)),
	isEmptyValue(reflect.ValueOf("h")),
	isEmptyValue(reflect.ValueOf(true)),
	isEmptyValue(reflect.ValueOf(0.1)),
	isEmptyValue(reflect.ValueOf(time.Now())),
}

func TestIsEmptyValue(t *testing.T) {
	for i, isEmpty := range emptyValues {
		if !isEmpty {
			t.Errorf("Value %d should be empty", i)
		}
	}
	for i, isNotEmpty := range nonEmptyValues {
		if isNotEmpty {
			t.Errorf("Value %d should not be empty", i)
		}
	}
}
