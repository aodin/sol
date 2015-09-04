package sol

import (
	"fmt"
)

// TODO What are the actual rules for column and table names?

func isValidColumnName(name string) error {
	if name == "" {
		return fmt.Errorf("sol: column names cannot be blank")
	}
	return nil
}

func isValidTableName(name string) error {
	if name == "" {
		return fmt.Errorf("sol: table names cannot be blank")
	}
	return nil
}
