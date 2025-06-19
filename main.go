package main

import (
	"context"
	"message-server/internal/controller/router"
	"message-server/internal/repository"
	"message-server/internal/usecases"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("DB_URL not set in .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic(err)
	}

	roomRepository := repository.NewRoomRepository(pool)
	authRepository := repository.NewAuthRepository(pool)
	listingRepository := repository.NewListingRepository(pool)
	fileRepository := repository.NewFileRepository(&repository.FileRepositoryConfig{
		AccountID: os.Getenv("R2_ACCOUNT_ID"),
		AccessKey: os.Getenv("R2_ACCESS_KEY"),
		SecretKey: os.Getenv("R2_SECRET_KEY"),
		Bucket:    os.Getenv("R2_BUCKET_NAME"),
	})
	userRepository := repository.NewUserRepository(pool)

	roomUseCase := usecases.NewRoomUseCase(roomRepository, authRepository, listingRepository)
	authUseCase := usecases.NewAuthUseCase(authRepository)
	listingUseCase := usecases.NewListingUseCase(listingRepository, fileRepository)
	fileUseCase := usecases.NewFileUseCase(fileRepository)
	userUseCase := usecases.NewUserUseCase(userRepository, fileRepository, authRepository)

	router := router.NewRouter(roomUseCase, authUseCase, listingUseCase, fileUseCase, userUseCase)
	router.Run(":" + port)
}
