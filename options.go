package sol

import "strings"

// Common options
const (
	OmitEmpty = "omitempty"
)

type options []string

func (o options) Has(option string) bool {
	for _, opt := range o {
		if opt == option {
			return true
		}
	}
	return false
}

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
