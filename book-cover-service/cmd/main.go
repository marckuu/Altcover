package main

import (
	"book-cover-service/core/db"
	"book-cover-service/core/logger"
	"book-cover-service/core/shared"
	"book-cover-service/core/shared/messaging"
	"book-cover-service/internal/auth/repositories"
	repositories2 "book-cover-service/internal/books/repositories"
	services2 "book-cover-service/internal/books/services"
	repositories3 "book-cover-service/internal/covers/repositories"
	services3 "book-cover-service/internal/covers/services"
	repositories4 "book-cover-service/internal/reactions/repositories"
	services4 "book-cover-service/internal/reactions/services"
	repositories5 "book-cover-service/internal/snapshots/repositories"
	"book-cover-service/internal/snapshots/services"
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

	coverRepository := repositories3.NewCoverRepository(conn)
	coverLikeRepository := repositories4.NewCoverLikeRepository(conn)
	favoritesRepository := repositories4.NewFavoritesRepository(conn)
	designerProfileSnapshotRepository := repositories5.NewDesignerProfileSnapshotRepository(conn)
	bookRepository := repositories2.NewBookRepository(conn)

	serverManager := NewServerManager(
		services3.NewCoverService(coverRepository),
		services4.NewCoverLikeService(coverLikeRepository),
		services4.NewFavoriteService(favoritesRepository, coverRepository),
		services.NewDesignerProfileSnapshotService(designerProfileSnapshotRepository),
		services2.NewBookService(bookRepository),
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
