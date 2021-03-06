package sol

import (
	"fmt"
	"log"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

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

var _ types.Type = FKElem{}

// FKElem is an internal type representation. It implements the
// types.Type interface so it can be used in CREATE TABLE statements.
type FKElem struct {
	name       string
	col        ColumnElem
	datatype   types.Type
	table      *TableElem // the parent table of the key
	references *TableElem // the table the key references
	onDelete   *fkAction
	onUpdate   *fkAction
}

// Create returns the element's syntax for a CREATE TABLE statement.
func (fk FKElem) Create(d dialect.Dialect) (string, error) {
	// Compile the type
	ct, err := fk.datatype.Create(d)
	if err != nil {
		return "", err
	}
	compiled := fmt.Sprintf(
		`%s %s REFERENCES %s(%s)`,
		fk.name,
		ct,
		fk.col.Table().Name(),
		fk.col.Name(),
	)
	if fk.onDelete != nil {
		compiled += fmt.Sprintf(" ON DELETE %s", *fk.onDelete)
	}
	if fk.onUpdate != nil {
		compiled += fmt.Sprintf(" ON UPDATE %s", *fk.onUpdate)
	}
	return compiled, nil
}

func (fk FKElem) ForeignName() string {
	return fk.col.Name()
}

// References returns the table that this foreign key references.
func (fk FKElem) References() *TableElem {
	return fk.references
}

// Modify implements the TableModifier interface. It creates a column and
// adds the same column to the create array.
func (fk FKElem) Modify(tabular Tabular) error {
	if tabular == nil || tabular.Table() == nil {
		return fmt.Errorf("sol: foreign keys cannot modify a nil table")
	}
	table := tabular.Table() // Get the dialect neutral table

	if err := isValidColumnName(fk.name); err != nil {
		return err
	}

	// Add the table to the foreign key
	if fk.table != nil && fk.table != table {
		return fmt.Errorf(
			"sol: foreign key %s already belongs to table %s",
			fk.name, fk.table.name,
		)
	}
	fk.table = table

	// Create the column for this table
	col := ColumnElem{
		name:     fk.name,
		table:    table,
		datatype: fk.datatype,
	}

	// Add the column to the table
	var err error
	if table.columns, err = table.columns.Add(col); err != nil {
		return err
	}

	// Add the type to the table creates
	table.creates = append(table.creates, fk)

	// Add it to the list of foreign keys
	table.fks = append(table.fks, fk)

	// Add the current table to the references of the foreign key table
	fk.references.referencedBy = append(fk.references.referencedBy, fk)

	return nil
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

// ForeignKey creates a FKElem from the given name and column/table.
// If given a column, it must already have its table assigned.
// If given a table, it must have one and only one primary key
func ForeignKey(name string, fk Selectable, datatypes ...types.Type) FKElem {
	if fk == nil {
		log.Panic("sol: inline foreign key was given a nil Selectable")
	}

	var col ColumnElem
	switch t := fk.(type) {
	case Columnar:
		col = t.Column()
	case Tabular:
		// If the given Selectable was a Table - use its primary key
		// TODO Can foreign keys reference multiple columns?
		pk := t.Table().PrimaryKey()
		if len(pk) != 1 {
			log.Panic(
				"sol: a table used directly in ForeignKey must have one and only one primary key column",
			)
		}
		col = t.Table().C(pk[0])
	default:
		log.Panicf("sol: unknown Selectable %T used in ForeignKey", t)
	}

	if col.Table() == nil {
		log.Panic(
			"sol: a column must have a table before being used in ForeignKey",
		)
	}

	// Allow an overriding datatype
	datatype := col.Type()
	if len(datatypes) > 0 {
		// TODO what should happen if multiple datatypes are given?
		datatype = datatypes[0]
	}

	return FKElem{
		name:       name,
		col:        col,
		datatype:   datatype,
		references: col.Table(),
	}
}

// SelfFKElem allows a table to have a foreign key to itself. willReference
// is a placeholder for the column the self-referential foreign key
// will reference
type SelfFKElem struct {
	FKElem
	willReference string
}

// Modify implements the TableModifier interface. It creates a column and
// adds the same column to the create array, will adding the referencing
// table and column
func (fk SelfFKElem) Modify(tabular Tabular) error {
	if tabular == nil || tabular.Table() == nil {
		return fmt.Errorf("sol: self foreign keys cannot modify a nil table")
	}
	table := tabular.Table() // Get the dialect neutral table

	fk.FKElem.table = table
	fk.FKElem.references = table

	// Does the reference column exist?
	fk.FKElem.col = table.C(fk.willReference)
	if fk.FKElem.col.IsInvalid() {
		return fmt.Errorf("sol: no column %s exists on table %s - is it created after the foreign key?", fk.willReference, table.Name())
	}

	// Set the datatype to the referenced column datatype - unless it has
	// already been set during construction
	if fk.FKElem.datatype == nil {
		fk.FKElem.datatype = fk.FKElem.col.Type()
	}

	return fk.FKElem.Modify(table)
}

// SelfForeignKey creates a self-referential foreign key
func SelfForeignKey(name, ref string, datatypes ...types.Type) SelfFKElem {
	// Allow the type to be overridden by a single optional type
	var datatype types.Type
	if len(datatypes) > 0 {
		datatype = datatypes[0]
	}

	return SelfFKElem{
		FKElem: FKElem{
			name:     name,
			datatype: datatype,
			// col and references will be added during Modify
		},
		willReference: ref,
	}
}
