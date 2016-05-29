package mysql

import (
	"os"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aodin/sol"
)

const travisCI = "root@tcp(127.0.0.1:3306)/sol_test?parseTime=true"

var conn *sol.DB
var once sync.Once

// getConn returns a MySQL connection pool
func getConn(t *testing.T) *sol.DB {
	// Check if an ENV VAR has been set, otherwise, use travis
	credentials := os.Getenv("SOL_TEST_MYSQL")
	if credentials == "" {
		credentials = travisCI
	}

	once.Do(func() {
		var err error
		if conn, err = sol.Open("mysql", credentials); err != nil {
			t.Fatalf("Failed to open connection: %s", err)
		}
		conn.SetMaxOpenConns(20)
	})
	return conn
}

func TestMySQL(t *testing.T) {
	conn := getConn(t)
	defer conn.Close()
	sol.IntegrationTest(t, conn)
}
