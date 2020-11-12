package main

import (
	"alchemy/galacticdb/db"
	"log"
	"alchemy/galacticdb/server"
	"net/http"
)

func main() {
	var port = "8081"
	s := db.GetSQLStore()
	r := server.GetRouter(s)
	log.Printf("Running golang server on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}