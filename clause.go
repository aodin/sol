package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

type Clause interface {
	Compiles
}

// ArrayClause is any number of clauses with a column join
type ArrayClause struct {
	clauses []Clause
	sep     string
}

func (c ArrayClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c ArrayClause) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	compiled := make([]string, len(c.clauses))
	var err error
	for i, clause := range c.clauses {
		compiled[i], err = clause.Compile(d, ps)
		if err != nil {
			return "", err
		}
	}
	return strings.Join(compiled, c.sep), nil
}

// AllOf joins the given clauses with 'AND' and wraps them in parentheses
func AllOf(clauses ...Clause) Clause {
	return FuncClause{Inner: ArrayClause{clauses, " AND "}}
}

// AnyOf joins the given clauses with 'OR' and wraps them in parentheses
func AnyOf(clauses ...Clause) Clause {
	return FuncClause{Inner: ArrayClause{clauses, " OR "}}
}

// BinaryClause is two clauses with a separator
type BinaryClause struct {
	Pre, Post Clause
	Sep       string
}

func (c BinaryClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

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

func (c FuncClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c FuncClause) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	cc, err := c.Inner.Compile(d, ps)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", c.Name, cc), nil
}

type UnaryClause struct {
	Pre Clause
	Sep string
}

func (c UnaryClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

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
