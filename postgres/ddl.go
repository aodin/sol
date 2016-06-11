package postgres

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/aodin/sol"
	"github.com/aodin/sol/types"
)

var ExcludeSchema = []string{"pg_catalog", "information_schema"}

var (
	cardinal   = types.New("information_schema.cardinal_number")
	identifier = types.New("information_schema.sql_identifier")
	chardata   = types.New("information_schema.character_data")
	yesorno    = types.New("information_schema.yes_or_no")
)

var columnsInfo = Table("information_schema.columns",
	// TODO Currently only a subset of all columns
	sol.Column("table_schema", identifier),
	sol.Column("table_name", identifier),
	sol.Column("column_name", identifier),
	sol.Column("ordinal_position", cardinal),
	sol.Column("column_default", chardata),
	sol.Column("is_nullable", yesorno),
	sol.Column("data_type", chardata),
	sol.Column("character_maximum_length", cardinal),
	sol.Column("character_octet_length", cardinal),
	sol.Column("numeric_precision", cardinal),
	sol.Column("numeric_precision_radix", cardinal),
	sol.Column("numeric_scale", cardinal),
	sol.Column("datetime_precision", cardinal),
)

type YesOrNo bool

func (field *YesOrNo) Scan(value interface{}) error {
	// TODO attempt parsing as a standard bool
	switch string(value.([]uint8)) {
	case "YES":
		*field = true
	case "NO":
		*field = false
	default:
		return fmt.Errorf("postgres: unable to match %s as YES or NO", value)
	}
	return nil
}

func (field YesOrNo) Value() (driver.Value, error) {
	if field {
		return "YES", nil
	}
	return "NO", nil
}

type ddl struct {
	TableSchema            string
	TableName              string
	ColumnName             string
	OrdinalPosition        int
	ColumnDefault          sql.NullString
	IsNullable             YesOrNo // returns as YES / NO
	DataType               string
	CharacterMaximumLength sql.NullInt64
	CharacterOctetLength   sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericPrecisionRadix  sql.NullInt64
	NumericScale           sql.NullInt64
	DatetimePrecision      sql.NullInt64
}

// ParseDDL will parse the DDL of the given connection. By default, it
// will parse all schemas except for 'pg_catalog' and 'information_schema'.
// If any schemas are given, it will parse only those.
func ParseDDL(conn sol.Conn, schemas ...string) ([]*TableElem, error) {
	stmt := columnsInfo.Select().OrderBy(
		columnsInfo.C("table_name"),
		columnsInfo.C("ordinal_position"),
	)
	if len(schemas) == 0 {
		stmt = stmt.Where(columnsInfo.C("table_schema").NotIn(ExcludeSchema))
	} else {
		stmt = stmt.Where(columnsInfo.C("table_schema").In(schemas))
	}

	var ddls []ddl
	if err := conn.Query(stmt, &ddls); err != nil {
		return nil, err
	}

	// Group the columns by table and schema
	type tableDDL struct {
		Schema, Name string
	}

	elements := make(map[tableDDL][]ddl)

	for _, ddl := range ddls {
		key := tableDDL{Schema: ddl.TableSchema, Name: ddl.TableName}
		elements[key] = append(elements[key], ddl)
	}

	// Parse each table
	var tables []*TableElem
	for tableDDL, ddls := range elements {
		var modifiers []sol.Modifier
		for _, ddl := range ddls {
			var final types.Type

			// Create a new base type for each ddl
			base := types.New(ddl.DataType)
			if !ddl.IsNullable {
				base.SetNotNull()
			}
			if ddl.ColumnDefault.Valid {
				base.SetDefault(ddl.ColumnDefault.String)
			}

			// Perform postgres specific type conversions (future use)
			switch ddl.DataType {
			case "boolean":
				final = types.BooleanType{BaseType: base}
			case "character varying", "text":
				final = types.CharacterType{BaseType: base}
			case "integer": // TODO way more numeric types
				final = types.NumericType{BaseType: base}
			case "date", "datetime", "timestamp without time zone":
				final = timestamp{BaseType: base}
			case "timestamp with time zone":
				base.SetName("timestamp")
				final = (timestamp{BaseType: base}).WithTimezone()
			case "json":
				final = json{BaseType: base}
			case "uuid":
				final = uuid{BaseType: base}
			default:
				final = base

				// TODO serial types and many more
			}

			// Create the column - TODO and any other modifiers
			modifiers = append(modifiers, sol.Column(ddl.ColumnName, final))
		}
		if len(modifiers) != 0 {
			tables = append(tables, Table(tableDDL.Name, modifiers...))
		}
	}

	return tables, nil
}
