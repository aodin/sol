package sol

import (
	"sort"

	"github.com/stretchr/testify/assert"
)

// Values is a map of column names to parameters.
type Values map[string]interface{}

// Diff returns the values in v that differ from the values in other.
// ISO 31-11: v \ other
func (v Values) Diff(other Values) Values {
	diff := Values{}
	for key, value := range v {
		if !assert.ObjectsAreEqual(value, other[key]) {
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
