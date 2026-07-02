package main

import (
	"auth-service/db"
	repositories2 "auth-service/db/repositories"
	"auth-service/server"
	"auth-service/services"
	"auth-service/shared/messaging"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	addresses := []string{os.Getenv("KAFKA_BROKER1_URL"), os.Getenv("KAFKA_BROKER2_URL"), os.Getenv("KAFKA_BROKER3_URL")}
	topic := os.Getenv("PROFILES_TOPIC_NAME")

	producer, err := messaging.NewProducer(addresses, topic)
	if err != nil {
		fmt.Printf("Не удалось создать продюсера: %v", err)
		return
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	conn, err := db.CreateConnection(ctx)
	if err != nil {
		fmt.Printf("Не удалось создать подключение к бд: %v", err)
		return
	}

	serverManager := server.NewServerManager(
		services.NewTokenService(repositories2.NewTokenRepository(conn), producer),
		services.NewUserService(repositories2.NewUserRepository(conn), producer),
		services.NewDesignerProfileService(repositories2.NewDesignerProfileRepository(conn), producer),
		repositories2.NewJWTManager(),
		ctx,
	)

	serverManager.StartServer()
}
