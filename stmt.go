package sol

import "fmt"

// Stmt is the base of all statements, including SELECT, UPDATE, DELETE, and
// INSERT statements
type Stmt struct {
	errs stmtErrors
}

// AddMeta adds a meta errors to the Stmt errors
func (stmt *Stmt) AddMeta(msg string, args ...interface{}) {
	// Create errs if they don't exist
	if stmt.errs.fields == nil {
		stmt.errs = stmtErrors{fields: make(map[fieldError]string)}
	}
	stmt.errs.meta = append(stmt.errs.meta, fmt.Sprintf(msg, args...))
}

// Error returns the statement's inner error
func (stmt Stmt) Error() error {
	if stmt.errs.Exist() {
		return stmt.errs
	}
	return nil
}

// TODO error setter

// ConditionalStmt includes SELECT, DELETE, and UPDATE statements
type ConditionalStmt struct {
	Stmt
	where Clause
}

// Conditional returns the statement's conditional Clause
func (stmt ConditionalStmt) Conditional() Clause {
	return stmt.where
}
