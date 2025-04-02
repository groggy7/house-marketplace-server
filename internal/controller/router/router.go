package router

import (
	"message-server/internal/controller"
	"message-server/internal/controller/auth"
	"message-server/internal/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(roomUseCase *usecases.RoomUseCase, authUseCase *usecases.AuthUseCase) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	wsServer := controller.InitMessageServer(roomUseCase, authUseCase)
	roomServer := controller.InitRoomServer(roomUseCase)
	authHandler := controller.NewAuthHandler(authUseCase)

	public := router.Group("")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/ws", wsServer.StartWebSocketServer)
	}

	protected := router.Group("")
	protected.Use(auth.JWTAuthMiddleware())
	{
		protected.POST("/room", roomServer.CreateRoom)
		protected.GET("/room/:customer_id", roomServer.GetRooms)
		protected.GET("/room/messages/:room_id", roomServer.GetRoomMessages)
	}

	return router
}
