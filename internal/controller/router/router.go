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
	fileUseCase *usecases.FileUseCase,
	userUseCase *usecases.UserUseCase,
) *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://house-marketplace-581ed5aac951.herokuapp.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	wsHandler := controller.InitMessageHandler(roomUseCase, authUseCase)
	roomHandler := controller.InitRoomHandler(roomUseCase)
	authHandler := controller.NewAuthHandler(authUseCase)
	listingHandler := controller.NewListingHandler(listingUseCase)
	fileHandler := controller.NewFileHandler(fileUseCase)
	userHandler := controller.NewUserHandler(userUseCase)

	public := router.Group("")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/ws", wsHandler.StartWebSocketServer)

		public.GET("/listing", listingHandler.GetListings)
		public.GET("/listing/:id", listingHandler.GetListingByID)
	}

	protected := router.Group("")
	protected.Use(auth.JWTAuthMiddleware())
	{
		protected.GET("/user", authHandler.CheckIsLoggedIn)
		protected.PUT("/user/info", userHandler.UpdateUserInfo)
		protected.PUT("/user/avatar", userHandler.UpdateUserAvatar)

		protected.POST("/logout", authHandler.Logout)
		protected.POST("/room", roomHandler.CreateRoom)
		protected.GET("/room", roomHandler.GetRooms)
		protected.GET("/room/messages/:room_id", roomHandler.GetRoomMessages)

		protected.POST("/listing", listingHandler.CreateListing)
		protected.PUT("/listing/:id", listingHandler.UpdateListing)
		protected.DELETE("/listing/:id", listingHandler.DeleteListing)

		protected.POST("/bookmark/:listing_id", listingHandler.BookmarkListing)
		protected.DELETE("/bookmark/:listing_id", listingHandler.UnbookmarkListing)
		protected.GET("/bookmark", listingHandler.GetBookmarkedListings)

		protected.POST("/file", fileHandler.UploadListingPicture)
		protected.DELETE("/file", fileHandler.DeleteFile)
	}

	return router
}
