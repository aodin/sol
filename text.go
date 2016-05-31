package sol

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/aodin/sol/dialect"
)

// paramsRegex will match any words with a leading colon.
// Since Go's re2 has no lookbehind / lookahead assertions, we'll match
// any number of leading colons - which will include type casts -
// and filter afterwards
// TODO this regular expression should be replaced by a parser
var paramsRegex = regexp.MustCompile(`(:)+(\w+)`)

// TextStmt allows the creation of custom SQL statements
type TextStmt struct {
	Stmt
	text   string
	values Values
}

// String outputs the parameter-less statement in a neutral dialect.
func (stmt TextStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile outputs the statement using the given dialect and parameters.
func (stmt TextStmt) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	// Select the parameters from the statement and replace them
	// with dialect specific parameters
	// Also alias the values (TODO hacky and probably too friendly)
	aliases := Values{}
	for key, value := range stmt.values {
		aliases[camelToSnake(key)] = value
	}

	replacer := func(match string) string {
		// Remove any matches with more than one leading colon
		if strings.LastIndex(match, ":") != 0 {
			return match
		}

		key := match[1:]

		// Parameter names must match value keys exactly or with
		// camel to snake conversion
		value, exists := stmt.values[key]
		if !exists {
			if value, exists = aliases[key]; !exists {
				stmt.AddMeta("sol: missing value for parameter '%s'", key)
			}
		}

		param := &Parameter{Value: value}
		replacement, err := param.Compile(d, ps)
		if err != nil {
			stmt.AddMeta(err.Error())
		}
		return replacement
	}
	compiled := paramsRegex.ReplaceAllStringFunc(stmt.text, replacer)
	return compiled, stmt.Error()
}

// Values sets the values of the statement. They can be given as either
// Values or struct types
func (stmt TextStmt) Values(obj interface{}) TextStmt {
	elem := reflect.Indirect(reflect.ValueOf(obj))

	// Examine allowed types
	var unsupported bool
	switch elem.Kind() {
	case reflect.Map:
		switch converted := obj.(type) {
		case Values:
			stmt.values = converted
		case *Values:
			stmt.values = *converted
		default:
			unsupported = true
		}
	case reflect.Struct:
		var err error
		if stmt.values, err = ValuesOf(obj); err != nil {
			stmt.AddMeta(err.Error())
			return stmt
		}
	default:
		unsupported = true
	}

	if unsupported {
		stmt.AddMeta(
			"sol: unsupported type %T for inserted values - accepted types: struct types or Values",
			obj,
		)
	}
	return stmt
}

// Text creates a TextStmt with custom SQL
func Text(text string, values ...Values) TextStmt {
	merged := Values{}
	for _, val := range values {
		merged = merged.Merge(val)
	}
	return TextStmt{text: text, values: merged}
}
