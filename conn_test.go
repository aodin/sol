package sol

import (
	"database/sql"
	"testing"
)

var _ executer = &sql.DB{}
var _ executer = &sql.Tx{}

func TestConn(t *testing.T) {}
