package config

import "fmt"

// MySQLConfig schema
type MySQLConfig struct {
	Host     string `name:"host" help:"Mysql host" env:"MYSQL_HORT"`
	Database string `name:"database" help:"Mysql database" env:"MYSQL_DATABASE"`
	Port     int    `name:"port" help:"Mysql port" env:"MYSQL_PORT"`
	Username string `name:"username" help:"Mysql username" env:"MYSQL_USERNAME"`
	Password string `name:"password" help:"Mysql password" env:"MYSQL_PASSWORD"`
	Options  string `name:"options" help:"Mysql options" env:"MYSQL_OPTIONS" default:"7845"`
}

// String return MySQL connection url
func (m MySQLConfig) String() string {
	return fmt.Sprintf("mysql://%s", m.DSN())
}

// DSN return Data Source Name
func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", m.Username, m.Password, m.Host, m.Port, m.Database, m.Options)
}
