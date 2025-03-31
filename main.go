package main

import (
	"message-server/db"
	"message-server/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewChatDB()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	wsServer := server.InitMessageServer(db)
	roomServer := server.InitRoomServer(db)

	router.GET("/ws", wsServer.StartWebSocketServer)

	router.POST("/create_room", roomServer.CreateRoom)
	router.POST("/get_rooms", roomServer.GetRooms)
	router.GET("/messages", roomServer.GetRoomMessages)

	router.Run(":80")
}
