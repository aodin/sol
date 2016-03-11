package sol

import (
	"encoding/json"
	"reflect"
	"sort"
)

// Values is a map of column names to parameters.
type Values map[string]interface{}

// Diff returns the values in v that differ from the values in other.
// ISO 31-11: v \ other
func (v Values) Diff(other Values) Values {
	diff := Values{}
	for key, value := range v {
		if !ObjectsAreEqual(value, other[key]) {
			diff[key] = value
		}
	}
	return diff
}

// Exclude removes the given keys and returns the remaining Values
func (v Values) Exclude(keys ...string) Values {
	safe := Values{}
ValueLoop:
	for key, value := range v {
		for _, k := range keys {
			if k == key {
				continue ValueLoop
			}
		}
		safe[key] = value
	}
	return safe
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

// Merge combines the Values{} types, with precedence to the Values given as
// a parameter
func (v Values) Merge(other Values) Values {
	merged := Values{}
	for key, value := range v {
		merged[key] = value
	}
	for key, value := range other {
		merged[key] = value
	}
	return merged
}

// MarshalJSON converts Values to JSON after converting all byte slices to
// a string type.
// By default, byte slices are JSON unmarshaled as base64.
// This is an issue since the postgres driver will scan string/varchar
// types as byte slices. Since Values{} should rarely be used within
// Go code, we're only modifying the JSON marshaler.
func (v Values) MarshalJSON() ([]byte, error) {
	for key, value := range v {
		if val, ok := value.([]byte); ok {
			v[key] = string(val)
		}
	}

	// Convert to prevent recursive marshaling
	return json.Marshal(map[string]interface{}(v))
}

// Values converts the given object to a Values{} type
func ValuesOf(obj interface{}) Values {
	values := Values{}

	elem := reflect.Indirect(reflect.ValueOf(obj))
	switch elem.Kind() {
	case reflect.Struct:
		fields := SelectFieldsFromElem(elem.Type())
		// TODO how to convert to db column name? Show Values even care?
		for _, field := range fields {
			var fieldElem reflect.Value = elem
			for _, name := range field.names {
				fieldElem = fieldElem.FieldByName(name)
			}
			// TODO Skip empty if omit empty...?
			values[field.column] = fieldElem.Interface()
		}
	case reflect.Map:
		// TODO Convert to Values - generalized map iteration?
	default:
		// TODO Return an error, panic, or silent?
	}
	return values
}
