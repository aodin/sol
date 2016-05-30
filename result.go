package sol

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Scanner is used for building mock result rows for testing
type Scanner interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...interface{}) error
}

// Result is returned by a database query - it embeds a Scanner
type Result struct {
	stmt string
	Scanner
}

// One returns a single row from Result. The destination must be a pointer
// to a struct or Values type. If there is only a single result column,
// then the destination can be a single native type.
func (r Result) One(obj interface{}) error {
	// Confirm that there is at least one row to return
	if ok := r.Next(); !ok {
		return sql.ErrNoRows
	}

	columns, err := r.Columns()
	if err != nil {
		return fmt.Errorf("sol: error returning columns from result: %s", err)
	}

	// Since maps are already pointers, they can be used as destinations
	// no matter what - as long as they are of type Values
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Map:
		if value.IsNil() {
			// TODO this might be a lie...
			return fmt.Errorf("sol: map types must be initialized before being used as destinations")
		}
		values, ok := obj.(Values)
		if !ok {
			return fmt.Errorf("sol: map types can be destinations only if they are of type Values")
		}

		// TODO scan directly into values?
		addr := make([]interface{}, len(columns))
		dest := make([]interface{}, len(columns))
		for i := range addr {
			dest[i] = &addr[i]
		}

		if err := r.Scan(dest...); err != nil {
			return fmt.Errorf("sol: error while scanning map: %s", err)
		}

		for i, name := range columns {
			values[name] = addr[i]
		}
		return r.Err()
	case reflect.Ptr: // Do nothing here
	default:
		return fmt.Errorf(
			"sol: received a non-pointer destination for Result.One",
		)
	}

	// Other types must be given as pointers to be valid destinations
	elem := reflect.Indirect(value)
	switch elem.Kind() {
	case reflect.Struct:
		fields := DeepFields(obj)

		// Create an interface pointer for each column's destination.
		// Unmatched values will be discarded
		dest := make([]interface{}, len(columns))
		var atLeastOneMatch bool
		for i, column := range columns {
			for _, field := range fields {
				// Match names either exactly and using camel to snake
				if field.Name == column || camelToSnake(field.Name) == column {
					dest[i] = field.Value.Addr().Interface()
					atLeastOneMatch = true
					break
				}
			}
		}

		// If nothing matched and the number of fields equals the number
		// columns, then blindly align columns and fields
		// TODO This may be too friendly of a feature
		if !atLeastOneMatch && len(fields) == len(columns) {
			for i, field := range fields {
				dest[i] = field.Value.Addr().Interface()
			}
		}

		if err := r.Scan(dest...); err != nil {
			return fmt.Errorf("sol: error while scanning struct: %s", err)
		}

	case reflect.Slice:
		return fmt.Errorf("sol: Result.One cannot scan into slices")
	default:
		if len(columns) != 1 {
			return fmt.Errorf(
				"sol: unsupported type %T for multi-column Result.One", obj,
			)
		}
		return r.Scan(elem.Addr().Interface()) // Scan directly into the elem
	}
	return r.Err()
}

// All returns all result rows scanned into the given interface, which
// must be a pointer to a slice of either structs or values. If there
// is only a single result column, then the destination can be a
// slice of a single native type (e.g. []int).
func (r Result) All(arg interface{}) error {
	argVal := reflect.ValueOf(arg)
	if argVal.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"sol: received a non-pointer destination for result.All",
		)
	}

	argElem := argVal.Elem()
	if argElem.Kind() != reflect.Slice {
		return fmt.Errorf(
			"sol: received a non-slice destination for result.All",
		)
	}

	// Get the type of the slice element
	elem := argElem.Type().Elem()

	columns, err := r.Columns()
	if err != nil {
		return fmt.Errorf(
			"sol: error returning columns from result: %s",
			err,
		)
	}

	switch elem.Kind() {
	case reflect.Struct:

		// Build the fields of the given struct
		// TODO this operation could be cached
		fields := FieldsFromElem(elem)

		// Align the fields to the selected columns
		// This will discard unmatched fields
		// TODO struct mode? error if not all columns were matched?
		aligned := AlignColumns(columns, fields)

		// If the aligned struct is empty, fallback to matching the fields in
		// order, but only if the length of the columns equals the fields
		if aligned.Empty() && len(columns) == len(fields) {
			aligned = fields
		}

		// Is there an existing slice element for this result?
		n := argElem.Len()

		// The number of results that hve been scanned
		var scanned int

		for r.Next() {
			if scanned < n {
				// Scan into an existing element
				newElem := argElem.Index(scanned)

				// Get an interface for each field and save a pointer to it
				dest := make([]interface{}, len(aligned))
				for i, field := range aligned {
					// If the field does not exist, the value will be discarded
					if !field.Exists() {
						dest[i] = &dest[i]
						continue
					}

					// Recursively get an interface to the elem's fields
					var fieldElem reflect.Value = newElem
					for _, name := range field.names {
						fieldElem = fieldElem.FieldByName(name)
					}
					dest[i] = fieldElem.Addr().Interface()
				}

				if err := r.Scan(dest...); err != nil {
					return err
				}
			} else {
				// Create a new slice element
				newElem := reflect.New(elem).Elem()

				// Get an interface for each field and save a pointer to it
				dest := make([]interface{}, len(aligned))
				for i, field := range aligned {
					// If the field does not exist, the value will be discarded
					if !field.Exists() {
						dest[i] = &dest[i]
						continue
					}

					// Recursively get an interface to the elem's fields
					var fieldElem reflect.Value = newElem
					for _, name := range field.names {
						fieldElem = fieldElem.FieldByName(name)
					}
					dest[i] = fieldElem.Addr().Interface()
				}

				if err := r.Scan(dest...); err != nil {
					return err
				}
				argElem.Set(reflect.Append(argElem, newElem))
			}
			scanned += 1
		}

	case reflect.Map:
		_, ok := arg.(*[]Values)
		if !ok {
			return fmt.Errorf("sol: slices of maps are only allowed if they are of type sol.Values")
		}

		for r.Next() {
			values := Values{}

			// TODO scan directly into values?
			addr := make([]interface{}, len(columns))
			dest := make([]interface{}, len(columns))
			for i := range addr {
				dest[i] = &addr[i]
			}

			if err := r.Scan(dest...); err != nil {
				return fmt.Errorf("sol: error while scanning map: %s", err)
			}

			for i, name := range columns {
				values[name] = addr[i]
			}

			argElem.Set(reflect.Append(argElem, reflect.ValueOf(values)))
		}

	default:
		// Single column results can be scanned into native types
		if len(columns) != 1 {
			return fmt.Errorf(
				"sol: unsupported destination for multi-column result: %s",
				elem.Kind(),
			)
		}
		for r.Next() {
			newElem := reflect.New(elem).Elem()
			if err := r.Scan(newElem.Addr().Interface()); err != nil {
				return err
			}
			argElem.Set(reflect.Append(argElem, newElem))
		}
	}

	return r.Err()
}
