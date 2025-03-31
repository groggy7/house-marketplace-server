package main

import (
	"message-server/db"
	"message-server/server"
	"net/http"
	"os"
)

func main() {
	db, err := db.NewChatDB()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	wsServer := server.InitMessageServer(db)
	roomServer := server.InitRoomServer(db)

	http.HandleFunc("/ws", wsServer.StartWebSocketServer)
	http.HandleFunc("POST /create_room", roomServer.CreateRoom)
	http.HandleFunc("GET /get_rooms", roomServer.GetRooms)

	http.ListenAndServe(":"+port, nil)
}