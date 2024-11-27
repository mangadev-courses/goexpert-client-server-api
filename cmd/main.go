package main

import (
	"client-server-api/pkg/client"
	"client-server-api/pkg/db"
	"client-server-api/pkg/server"
)

func main() {
	db := db.DbClient()
	defer db.Close()

	go server.Start(db)
	client.FetchExchangeRate("http://localhost:8080/cotacao")
}
