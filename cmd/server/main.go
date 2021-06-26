package main

import (
	"coding-challenge-go/server"
	"coding-challenge-go/server/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	cfg := config.Load()

	server.Server(cfg)
}
