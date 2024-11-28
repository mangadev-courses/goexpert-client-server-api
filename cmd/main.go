package main

import (
	"client-server-api/pkg/client"
	"client-server-api/pkg/db"
	"client-server-api/pkg/server"
)

func main() {
	db := db.DbClient()
	defer db.Close()

	ready := make(chan struct{})
	server.Start(db, ready)
	<-ready

	err := client.FetchExchangeRate("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
}
