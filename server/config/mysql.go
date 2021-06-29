package config

import "fmt"

// MySQLConfig schema
type MySQLConfig struct {
	Host     string
	Database string
	Port     int
	Username string
	Password string
	Options  string
}

// String return MySQL connection url
func (m MySQLConfig) String() string {
	return fmt.Sprintf("mysql://%s", m.DSN())
}

// DSN return Data Source Name
func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", m.Username, m.Password, m.Host, m.Port, m.Database, m.Options)
}
