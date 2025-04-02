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
	roomUseCase := usecases.NewRoomUseCase(roomRepository)

	authRepository := repository.NewAuthRepository(pool)
	authUseCase := usecases.NewAuthUseCase(authRepository)

	router := router.NewRouter(roomUseCase, authUseCase)
	router.Run(":" + port)
}
