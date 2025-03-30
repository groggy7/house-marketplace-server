package main

import (
	"message-server/db"
	"message-server/server"
	"net/http"
)

func main() {
	pool, err := db.CreateDBConnection()
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	server := server.InitMessageServer(pool)
	http.HandleFunc("/ws", server.StartWebSocketServer)
	http.ListenAndServe(":8080", nil)
}