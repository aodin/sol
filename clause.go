package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// Clause is the interface that all structural components of a statement
// must implement
type Clause interface {
	Compiles
}

// String implements the Clause interface for strings
// TODO dialect specific, safe-escape
type String string

func (str String) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	return fmt.Sprintf("'%s'", str), nil
}

// ArrayClause is any number of clauses with a column join
type ArrayClause struct {
	clauses []Clause
	sep     string
	post    bool // Used for wrap
}

var _ Clause = ArrayClause{}

// String returns the parameter-less ArrayClause in a neutral dialect.
func (c ArrayClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile returns the ArrayClause as a compiled string using
// the given Dialect - possibly with an error. Any parameters will
// be appended to the given Parameters.
func (c ArrayClause) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	compiled := make([]string, len(c.clauses))
	var err error
	for i, clause := range c.clauses {
		if compiled[i], err = clause.Compile(d, ps); err != nil {
			return "", err
		}
	}
	return strings.Join(compiled, c.sep), nil
}

// Wrap implements the Operator interface
// TODO Should Wrap operate on a clause and return a clause?
func (c ArrayClause) Wrap(str string) string {
	// Compile all clauses
	clauses := make([]string, len(c.clauses))
	var err error
	for i, clause := range c.clauses {
		// TODO Is it dialect specific or not?
		if clauses[i], err = clause.Compile(nil, Params()); err != nil {
			return "" // silent error?
		}
	}
	if c.post {
		return strings.Join(append(clauses, str), c.sep)
	}
	return strings.Join(append([]string{str}, clauses...), c.sep)
}

// AllOf joins the given clauses with 'AND' and wraps them in parentheses
func AllOf(clauses ...Clause) Clause {
	return FuncClause{Inner: ArrayClause{clauses: clauses, sep: " AND "}}
}

// AnyOf joins the given clauses with 'OR' and wraps them in parentheses
func AnyOf(clauses ...Clause) Clause {
	return FuncClause{Inner: ArrayClause{clauses: clauses, sep: " OR "}}
}

// BinaryClause is two clauses with a separator
type BinaryClause struct {
	Pre, Post Clause
	Sep       string
}

var _ Clause = BinaryClause{}

// String returns the parameter-less BinaryClause in a neutral dialect.
func (c BinaryClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile returns the BinaryClause as a compiled string using
// the given Dialect - possibly with an error. Any parameters will
// be appended to the given Parameters.
func (c BinaryClause) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	var pre, post string
	var err error
	if c.Pre != nil {
		pre, err = c.Pre.Compile(d, ps)
		if err != nil {
			return "", err
		}
	}
	if c.Post != nil {
		post, err = c.Post.Compile(d, ps)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s%s%s", pre, c.Sep, post), nil
}

type FuncClause struct {
	Inner Clause
	Name  string
}

var _ Clause = FuncClause{}

// String returns the parameter-less FuncClause in a neutral dialect.
func (c FuncClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile returns the FuncClause as a compiled string using
// the given Dialect - possibly with an error. Any parameters will
// be appended to the given Parameters.
func (c FuncClause) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	cc, err := c.Inner.Compile(d, ps)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", c.Name, cc), nil
}

// Wrap implements the Operator interface
func (c FuncClause) Wrap(str string) string {
	return fmt.Sprintf("%s(%s)", c.Name, str)
}

type UnaryClause struct {
	Pre Clause
	Sep string
}

var _ Clause = UnaryClause{}

// String returns the parameter-less UnaryClause in a neutral dialect.
func (c UnaryClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile returns the UnaryClause as a compiled string using
// the given Dialect - possibly with an error. Any parameters will
// be appended to the given Parameters.
func (c UnaryClause) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	var pre string
	var err error
	if c.Pre != nil {
		pre, err = c.Pre.Compile(d, ps)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s%s", pre, c.Sep), nil
}
