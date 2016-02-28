package sol

type Comparable interface {
	Equals(interface{}) BinaryClause
	GTE(interface{}) BinaryClause
}
