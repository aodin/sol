package mysql

import (
	"os"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aodin/sol"
)

const travisCI = "root@tcp(127.0.0.1:3306)/sol_test?parseTime=true"

var testconn *sol.DB
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
		if testconn, err = sol.Open("mysql", credentials); err != nil {
			t.Fatalf("Failed to open connection: %s", err)
		}
		testconn.SetMaxOpenConns(20)
	})
	return testconn
}

// TestMySQL performs the standard integration test
func TestMySQL(t *testing.T) {
	conn := getConn(t)
	defer conn.Close()
	sol.IntegrationTest(t, conn, true)
}
