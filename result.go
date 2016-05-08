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

// Result is returned by a database query - it implements Scanner
type Result struct {
	stmt string
	Scanner
}

// One returns a single row from Result. The destination interface must be
// a pointer to a struct or a native type.
func (r *Result) One(arg interface{}) error {
	// Confirm that there is at least one row to return
	if ok := r.Next(); !ok {
		return sql.ErrNoRows
	}

	columns, err := r.Columns()
	if err != nil {
		return fmt.Errorf(
			"sol: error returning columns from result: %s",
			err,
		)
	}

	value := reflect.ValueOf(arg)
	if value.Kind() == reflect.Map {
		values, ok := arg.(Values)
		if !ok {
			return fmt.Errorf("sol: maps as destinations are only allowed if they are of type sol.Values")
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

	} else if value.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"sol: received a non-pointer destination for result.One",
		)
	}

	// Get the value of the given interface
	elem := reflect.Indirect(value)

	switch elem.Kind() {
	case reflect.Struct:
		// Build the fields of the given struct
		// TODO this operation could be cached
		fields := SelectFields(arg)

		// Align the fields to the selected columns
		// This will discard unmatched fields
		// TODO strict mode? error if not all columns were matched?
		aligned := AlignColumns(columns, fields)

		// If the aligned struct is empty, fallback to matching the fields in
		// order, but only if the length of the columns equals the fields
		if aligned.Empty() && len(columns) == len(fields) {
			aligned = fields
		}

		// Get an interface for each field and save a pointer to it
		dest := make([]interface{}, len(aligned))
		for i, field := range aligned {
			// If the field does not exist, the value will be discarded
			if !field.Exists() {
				dest[i] = &dest[i]
				continue
			}

			// Recursively get an interface to the elem's fields
			var fieldElem reflect.Value = elem
			for _, name := range field.names {
				fieldElem = fieldElem.FieldByName(name)
			}
			dest[i] = fieldElem.Addr().Interface()
		}

		if err := r.Scan(dest...); err != nil {
			return fmt.Errorf("sol: error while scanning struct: %s", err)
		}

	case reflect.Slice:
		return fmt.Errorf("sol: cannot scan single results into slices")

	default:
		if len(columns) != 1 {
			return fmt.Errorf(
				"sol: unsupported destination for multi-column result: %s",
				elem.Kind(),
			)
		}
		// Attempt to scan directly into the elem
		return r.Scan(elem.Addr().Interface())
	}
	return r.Err()
}

// All returns all result rows into the given interface, which must be a
// pointer to a slice of either structs, values, or a native type.
func (r *Result) All(arg interface{}) error {
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
		fields := SelectFieldsFromElem(elem)

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
