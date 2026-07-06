package main

import (
	"auth-service/core/db"
	logs "auth-service/core/logger"
	"auth-service/core/logger/zap"
	"auth-service/core/shared/messaging"
	"auth-service/internal/auth/repositories"
	"auth-service/internal/auth/repositories/interfaces/jwtCovers"
	services2 "auth-service/internal/auth/services"
	repositories2 "auth-service/internal/designerProfiles/repositories"
	"auth-service/internal/designerProfiles/services"
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

	logger, loggerCancel, err := logs.NewLogger("INFO")
	if err != nil {
		fmt.Printf("Не удалось создать логгер: %v", err)
		return
	}

	loggerCover := zap.NewZapLoggerCover(logger)

	defer func() {
		err = loggerCancel()
		if err != nil {
			logger.Error(fmt.Errorf("не удалось закрыть логгер: %w", err).Error())
			return
		}
		err = producer.Close()
		if err != nil {
			logger.Error(fmt.Errorf("не удалось закрыть продюсера: %w", err).Error())
			return
		}
	}()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	conn, err := db.CreateConnection(ctx)
	if err != nil {
		fmt.Printf("Не удалось создать подключение к бд: %v", err)
		return
	}

	serverManager := NewServerManager(
		services2.NewTokenService(repositories.NewTokenRepository(conn), producer),
		services2.NewUserService(repositories.NewUserRepository(conn), producer),
		services.NewDesignerProfileService(repositories2.NewDesignerProfileRepository(conn), producer),
		jwtCovers.NewJwtManagerCover(repositories.JwtManager{}),
		ctx,
		loggerCover,
	)

	serverManager.StartServer()
}
