package sol

import (
	"github.com/aodin/sol/dialect"
)

type Parameter struct {
	Value interface{}
}

// Parameter compilation is dialect dependent. For instance, dialects such
// as PostGres require the parameter index.
func (p *Parameter) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	ps.Add(p.Value)
	return d.Param(ps.Len() - 1), nil
}

type Parameters []interface{}

func (ps *Parameters) Add(param interface{}) {
	*ps = append(*ps, param)
}

func (ps *Parameters) Len() int {
	return len(*ps)
}

func Params() *Parameters {
	return &Parameters{}
}
