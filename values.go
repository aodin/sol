package sol

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/aodin/sol/dialect"
)

// Values is a map of column names to values. It can be used both as
// a source of values in INSERT, UPDATE, and Text statements, or as
// a destination for SELECT.
type Values map[string]interface{}

// Compile outputs the Values in a format for UPDATE using the given dialect
// and parameters.
func (v Values) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	keys := v.Keys()
	values := make([]string, len(keys))
	for i, key := range keys {
		param := NewParam(v[key])
		compiledParam, err := param.Compile(d, ps)
		if err != nil {
			return "", err
		}
		values[i] = fmt.Sprintf(`%s = %s`, key, compiledParam)
	}
	return strings.Join(values, ", "), nil
}

// Diff returns the values in v that differ from the values in other.
// ISO 31-11: v \ other
func (v Values) Diff(other Values) Values {
	diff := Values{}
	for key, value := range v {
		if !reflect.DeepEqual(value, other[key]) {
			diff[key] = value
		}
	}
	return diff
}

// Equals returns true if both the receiver and parameter Values have
// equal keys and values
func (v Values) Equals(other Values) bool {
	if len(v) != len(other) {
		return false
	}
	for key, a := range v {
		if b, ok := other[key]; !ok || !reflect.DeepEqual(a, b) {
			return false
		}
	}
	return true
}

// Exclude removes the given keys and returns the remaining Values
func (v Values) Exclude(keys ...string) Values {
	out := Values{}
ValueLoop:
	for key, value := range v {
		for _, k := range keys {
			if k == key {
				continue ValueLoop
			}
		}
		out[key] = value
	}
	return out
}

// Filters returns a Values type with key-values from the original Values
// that match the given keys
func (v Values) Filter(keys ...string) Values {
	out := Values{}
	for key, value := range v {
		for _, k := range keys {
			if k == key {
				out[key] = value
				break
			}
		}
	}
	return out
}

// Keys returns the keys of the Values map in alphabetical order.
func (v Values) Keys() []string {
	keys := make([]string, len(v))
	var i int
	for key := range v {
		keys[i] = key
		i += 1
	}
	sort.Strings(keys)
	return keys
}

// Merge combines the given Values without modifying the original Values.
// Precedence is given to the rightmost Values given as parameters.
func (v Values) Merge(others ...Values) Values {
	merged := Values{}
	for _, other := range append([]Values{v}, others...) {
		for key, value := range other {
			merged[key] = value
		}
	}
	return merged
}

// MarshalJSON converts Values to JSON after converting all byte slices to
// a string type. By default, byte slices are JSON unmarshaled as base64.
// This is an issue since the postgres driver will scan string/varchar
// types as byte slices - the current solution is to convert before output.
func (v Values) MarshalJSON() ([]byte, error) {
	for key, value := range v {
		if val, ok := value.([]byte); ok {
			v[key] = string(val)
		}
	}

	// Convert to prevent recursive marshaling
	return json.Marshal(map[string]interface{}(v))
}

// Reject is an alias for Exclude
func (v Values) Reject(keys ...string) Values {
	return v.Exclude(keys...)
}

// Values returns the values of the Values map in the alphabetical order
// of its keys.
func (v Values) Values() []interface{} {
	keys := v.Keys()
	values := make([]interface{}, len(keys))
	for i, key := range keys {
		values[i] = v[key]
	}
	return values
}

// Values converts the given object to a Values{} type
func ValuesOf(obj interface{}) (Values, error) {
	elem := reflect.Indirect(reflect.ValueOf(obj))
	switch elem.Kind() {
	case reflect.Map:
		// If the type is already Values, convert and return
		switch converted := obj.(type) {
		case Values:
			return converted, nil
		case *Values:
			return *converted, nil
		case map[string]interface{}:
			return Values(converted), nil
		case *map[string]interface{}:
			return Values(*converted), nil
		default:
			return nil, fmt.Errorf(
				"sol: unsupported map type %T for ValuesOf()", converted,
			)
		}
	case reflect.Struct:
		values := Values{}
		for _, field := range DeepFields(obj) {
			if field.IsOmittable() {
				continue // Skip zero values
			}
			values[field.Name] = field.Value.Interface()
		}
		return values, nil
	}
	return nil, fmt.Errorf("sol: unsupported type %T for ValuesOf()", obj)
}

// isEmptyValue is from Go's encoding/json package: encode.go
// Copyright 2010 The Go Authors. All rights reserved.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		t, ok := v.Interface().(time.Time)
		if ok {
			return t.IsZero()
		}
	}
	return false
}
