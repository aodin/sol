package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Database contains the fields needed to connect to a database.
type Database struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
}

// Credentials with return a string of credentials appropriate for Go's
// sql.Open function
func (db Database) Credentials() (string, string) {
	// Only add the key if there is a value
	var values []string
	if db.Host != "" {
		values = append(values, fmt.Sprintf("host=%s", db.Host))
	}
	if db.Port != 0 {
		values = append(values, fmt.Sprintf("port=%d", db.Port))
	}
	if db.Name != "" {
		values = append(values, fmt.Sprintf("dbname=%s", db.Name))
	}
	if db.User != "" {
		values = append(values, fmt.Sprintf("user=%s", db.User))
	}
	if db.Password != "" {
		values = append(values, fmt.Sprintf("password=%s", db.Password))
	}
	if db.SSLMode != "" {
		values = append(values, fmt.Sprintf("sslmode=%s", db.SSLMode))
	}
	return db.Driver, strings.Join(values, " ")
}

// Parse will create a Database using the file at the given path.
func Parse(filename string) (Database, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Database{}, err
	}
	return parse(f)
}

func parse(f io.Reader) (c Database, err error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &c)
	return
}
