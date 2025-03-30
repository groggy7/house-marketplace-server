package main

import (
	"context"
	"message-server/db"
	"message-server/server"
	"net/http"
)

func main() {
	conn := db.CreateDBConnection()
	defer conn.Close(context.Background())

	server := server.InitMessageServer(conn)
	http.HandleFunc("/ws", server.StartWebSocketServer)
	http.ListenAndServe(":8080", nil)
}