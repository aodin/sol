package sol

import ()

type fkAction string

// The following constants represent possible foreign key actions that can
// be used in ON DELETE and ON UPDATE clauses.
const (
	NoAction   fkAction = "NO ACTION"
	Restrict   fkAction = "RESTRICT"
	Cascade    fkAction = "CASCADE"
	SetNull    fkAction = "SET NULL"
	SetDefault fkAction = "SET DEFAULT"
)

// FKElem is an internal type representation. It implements the
// Creatable interface so it can be used in CREATE TABLE statements.
type FKElem struct {
	name     string
	col      *ColumnElem
	table    *TableElem // the parent table of the key
	refTable *TableElem // the table the key references
	onDelete *fkAction
	onUpdate *fkAction
}

// OnDelete adds an ON DELETE clause to the foreign key
func (fk FKElem) OnDelete(b fkAction) FKElem {
	fk.onDelete = &b
	return fk
}

// OnUpdate add an ON UPDATE clause to the foreign key
func (fk FKElem) OnUpdate(b fkAction) FKElem {
	fk.onUpdate = &b
	return fk
}
