package main

import (
	"message-server/db"
	"message-server/server"
	"net/http"
	"os"
)

func main() {
	pool, err := db.CreateDBConnection()
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	server := server.InitMessageServer(pool)
	http.HandleFunc("/ws", server.StartWebSocketServer)
	http.ListenAndServe(":"+port, nil)
}