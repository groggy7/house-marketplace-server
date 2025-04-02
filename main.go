package main

import (
	"context"
	"message-server/internal/controller"
	"message-server/internal/repository"
	"message-server/internal/room"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("DB_URL not set in .env")
	}

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	roomRepository := repository.NewRoomRepository(pool)
	roomService := room.NewRoomService(roomRepository)
	wsServer := controller.InitMessageServer(roomService)
	roomServer := controller.InitRoomServer(roomService)

	router.GET("/ws", wsServer.StartWebSocketServer)

	router.POST("/room", roomServer.CreateRoom)
	router.GET("/room/:customer_id", roomServer.GetRooms)
	router.GET("/room/messages/:room_id", roomServer.GetRoomMessages)

	router.Run(":80")
}
