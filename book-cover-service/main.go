package main

import (
	"book-cover-service/db"
	"book-cover-service/db/repositories"
	logs "book-cover-service/logger"
	"book-cover-service/server"
	"book-cover-service/services"
	"book-cover-service/shared"
	"book-cover-service/shared/messaging"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	addresses := []string{os.Getenv("KAFKA_BROKER1_URL"), os.Getenv("KAFKA_BROKER2_URL"), os.Getenv("KAFKA_BROKER3_URL")}
	topics := []string{os.Getenv("PROFILES_TOPIC_NAME")}
	consumerGroupID := os.Getenv("CONSUMER_GROUP_ID")

	logger, loggerCancel, err := logs.NewLogger("INFO")
	if err != nil {
		fmt.Printf("Не удалось создать логгер: %v", err)
		return
	}

	defer loggerCancel()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	conn, err := db.CreateConnection(ctx)
	if err != nil {
		fmt.Printf("Не удалось создать подключение к бд: %v", err)
		return
	}

	coverRepository := repositories.NewCoverRepository(conn)
	coverLikeRepository := repositories.NewCoverLikeRepository(conn)
	favoritesRepository := repositories.NewFavoritesRepository(conn)
	designerProfileSnapshotRepository := repositories.NewDesignerProfileSnapshotRepository(conn)
	bookRepository := repositories.NewBookRepository(conn)

	serverManager := server.NewServerManager(
		services.NewCoverService(coverRepository),
		services.NewCoverLikeService(coverLikeRepository),
		services.NewFavoriteService(favoritesRepository, coverRepository),
		services.NewDesignerProfileSnapshotService(designerProfileSnapshotRepository),
		services.NewBookService(bookRepository),
		repositories.NewJWTManager(),
		ctx,
		logger,
	)

	consumer, err := messaging.NewConsumer(addresses, consumerGroupID, topics)
	if err != nil {
		fmt.Printf("не удалось создать консюмера, %v", err)
		return
	}

	eventHandlers := shared.NewEventHandlers(ctx, designerProfileSnapshotRepository)

	// Добавить контекст, закрытие консюмера
	go func() {
		defer consumer.Close()
		if err = consumer.Consume(eventHandlers.HandleKafkaEvent); err != nil {
			return
			// Логировать ошибку
		}
	}()

	serverManager.StartServer()
}
