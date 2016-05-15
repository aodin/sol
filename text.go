package sol

import (
	"fmt"
	"regexp"

	"github.com/aodin/sol/dialect"
)

var paramsRegex = regexp.MustCompile(`:(\w+)`)

type TextStmt struct {
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
	// TODO aggregate errors?
	var err error

	// Select the parameters from the statement and replace them
	// with dialect specific parameters
	replacer := func(match string) string {
		// TODO the regex should ignore the colon
		key := match[1:]

		// Parameter names must match value keys
		value, exists := stmt.values[key]
		if !exists {
			err = fmt.Errorf("sol: missing value for parameter '%s'", key)
		}

		param := &Parameter{Value: value}
		compiled, paramErr := param.Compile(d, ps)
		if paramErr != nil {
			err = paramErr
		}
		return compiled
	}

	return paramsRegex.ReplaceAllStringFunc(stmt.text, replacer), err
}

func (stmt TextStmt) Values(values Values) TextStmt {
	stmt.values = values
	return stmt
}

func Text(text string, values ...Values) TextStmt {
	merged := Values{}
	for _, val := range values {
		merged = merged.Merge(val)
	}
	return TextStmt{text: text, values: merged}
}
