package sol

import ()

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
