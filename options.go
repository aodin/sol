package sol

import "strings"

// Common options
const (
	OmitEmpty  = "omitempty"  // Skip this field if it has a zero value
	OmitUpdate = "omitupdate" // Skip this field during updates
)

type options []string

// Equals tests equality - order matters
func (o options) Equals(other options) bool {
	if len(o) != len(other) {
		return false
	}
	for i, opt := range other {
		if opt != other[i] {
			return false
		}
	}
	return true
}

// Has returns true if the given options exists in the current options
func (o options) Has(option string) bool {
	for _, opt := range o {
		if opt == option {
			return true
		}
	}
	return false
}

// parseTag splits a DB struct tag into its name and options
func parseTag(tag string) (string, options) {
	parts := strings.Split(tag, ",")
	return parts[0], options(parts[1:])
}

// splitName separates the tag into table and column names.
// If no separator (.) is given, assume the name is column only.
func splitName(name string) (string, string) {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return "", parts[0]
}
