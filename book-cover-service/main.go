package main

import (
	"book-cover-service/db/repositories"
	"book-cover-service/server"
	"book-cover-service/services"
	"book-cover-service/shared"
	"book-cover-service/shared/messaging"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {

	addresses := []string{os.Getenv("KAFKA_BROKER1_URL"), os.Getenv("KAFKA_BROKER2_URL"), os.Getenv("KAFKA_BROKER3_URL")}
	topics := []string{os.Getenv("PROFILES_TOPIC_NAME")}
	consumerGroupID := os.Getenv("CONSUMER_GROUP_ID")

	ctx := context.Background()
	var conn *pgx.Conn

	coverRepository := repositories.NewCoverRepository(conn)
	coverLikeRepository := repositories.NewCoverLikeRepository(conn)
	favoritesRepository := repositories.NewFavoritesRepository(conn)

	serverManager := server.NewServerManager(
		services.NewCoverService(coverRepository),
		services.NewCoverLikeService(coverLikeRepository),
		services.NewFavoriteService(favoritesRepository, coverRepository),
		repositories.NewJWTManager(),
		ctx,
	)

	consumer, err := messaging.NewConsumer(addresses, consumerGroupID, topics)
	if err != nil {
		fmt.Printf("не удалось создать консюмера, %v", err)
		return
	}

	designerProfileRepository := repositories.NewDesignerProfileRepository()
	eventHandlers := shared.NewEventHandlers(designerProfileRepository)

	// Добавить контекст, закрытие консюмера и обработчик, который будет выполнять все нужные действия
	go func() {
		defer consumer.Close()
		if err = consumer.Consume(eventHandlers.HandleKafkaEvent); err != nil {
			return
			// Логировать ошибку
		}
	}()

	serverManager.StartServer()
}
