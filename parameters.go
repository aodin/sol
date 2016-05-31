package sol

import "github.com/aodin/sol/dialect"

// Parameter is a value that will be passed to the database using a
// dialect specific parameterization
type Parameter struct {
	Value interface{}
}

// Parameter compilation is dialect dependent. For instance, dialects such
// as PostGres require the parameter index.
func (p *Parameter) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	ps.Add(p.Value)
	return d.Param(ps.Len() - 1), nil
}

// NewParam creates a new *Parameter
func NewParam(value interface{}) *Parameter {
	return &Parameter{value}
}

// Parameters aggregated values before parameterization
type Parameters []interface{}

// Add adds a parameter
func (ps *Parameters) Add(param interface{}) {
	*ps = append(*ps, param)
}

// Len returns the length of the Parameters
func (ps *Parameters) Len() int {
	return len(*ps)
}

// Params creates a new Parameters
func Params() *Parameters {
	return &Parameters{}
}
