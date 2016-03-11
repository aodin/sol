package sol

import (
	"encoding/json"
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
