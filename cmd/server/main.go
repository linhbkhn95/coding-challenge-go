package main

import (
	_ "github.com/go-sql-driver/mysql"

	"coding-challenge-go/server"
	"coding-challenge-go/server/config"
)

func main() {

	cfg := config.Load()

	server.Server(cfg)
}
