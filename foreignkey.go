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
	col        Columnar
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
		`"%s" %s REFERENCES %s("%s")`,
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
func (fk FKElem) Modify(table *TableElem) error {
	if table == nil {
		return fmt.Errorf("sol: foreign keys cannot modify a nil table")
	}
	if err := isValidColumnName(fk.name); err != nil {
		return err
	}

	// Add the table to the foreign key
	if fk.table != nil {
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
	if err := table.columns.add(col); err != nil {
		return err
	}

	// Add the type to the table creates
	table.creates = append(table.creates, fk)

	// Add it to the list of foreign keys
	table.fks = append(table.fks, fk)

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
func ForeignKey(name string, fk Selectable, datatypes ...types.Type) FKElem {
	var col Columnar
	if fk == nil {
		log.Panic("sol: inline foreign key was given a nil Selectable")
	}
	columns := fk.Columns()
	if len(columns) == 0 {
		log.Panic(
			"sol: inline foreign key Selectable must have at least one column",
		)
	} else if len(columns) == 1 {
		col = columns[0]
		if col.Table() == nil {
			log.Panic(
				"sol: inline foreign key columns must have their table assigned before creation",
			)
		}
	} else {
		// Simply use the table of the first column
		// TODO This is a strange decision that will error silently
		table := columns[0].Table()
		pk := table.PrimaryKey()
		if len(pk) != 1 {
			log.Panic(
				"sol: inline foreign key tables must have one and only one primary key column",
			)
		}
		col = table.C(pk[0])
	}

	// Allow an overriding datatype
	datatype := col.Type()
	if len(datatypes) > 0 {
		datatype = datatypes[0]
	}

	return FKElem{
		name:       name,
		col:        col,
		datatype:   datatype,
		references: col.Table(),
	}
}
