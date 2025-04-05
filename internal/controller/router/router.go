package router

import (
	"message-server/internal/controller"
	"message-server/internal/controller/auth"
	"message-server/internal/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	roomUseCase *usecases.RoomUseCase,
	authUseCase *usecases.AuthUseCase,
	listingUseCase *usecases.ListingUseCase,
) *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://house-marketplace-581ed5aac951.herokuapp.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	wsServer := controller.InitMessageServer(roomUseCase, authUseCase)
	roomServer := controller.InitRoomServer(roomUseCase)
	authHandler := controller.NewAuthHandler(authUseCase)
	listingHandler := controller.NewListingHandler(listingUseCase)

	public := router.Group("")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/ws", wsServer.StartWebSocketServer)
	}

	protected := router.Group("")
	protected.Use(auth.JWTAuthMiddleware())
	{
		protected.GET("/user", authHandler.CheckIsLoggedIn)

		protected.POST("/logout", authHandler.Logout)
		protected.POST("/room", roomServer.CreateRoom)
		protected.GET("/room", roomServer.GetRooms)
		protected.GET("/room/messages/:room_id", roomServer.GetRoomMessages)

		protected.POST("/listing", listingHandler.CreateListing)
		protected.GET("/listing/:id", listingHandler.GetListingByID)
		protected.GET("/listing", listingHandler.GetListings)
		protected.PUT("/listing/:id", listingHandler.UpdateListing)
		protected.DELETE("/listing/:id", listingHandler.DeleteListing)

		protected.POST("/listing/:id/bookmark", listingHandler.BookmarkListing)
		protected.DELETE("/listing/:id/bookmark", listingHandler.UnbookmarkListing)
	}

	return router
}
