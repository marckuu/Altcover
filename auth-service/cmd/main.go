package main

import (
	"auth-service/core/db"
	logs "auth-service/core/logger"
	"auth-service/core/logger/zap"
	"auth-service/core/shared/messaging"
	repositories3 "auth-service/internal/admin/repositories"
	services3 "auth-service/internal/admin/services"
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

// @title Маркетплейс альтернативных обложек для книг
// @description Сервис аутентификации
// @version 1.0
// @contact.name markuu
// @contact.url https://github.com/marckuu

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
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
	coonPool, err := db.CreateConnPool(ctx, os.Getenv("DB_CONN_PATH"))
	if err != nil {
		fmt.Printf("Не удалось создать подключение к бд: %v", err)
		return
	}

	tokenService := services2.NewTokenService(repositories.NewTokenRepository(coonPool), producer)
	jwtManagerCover := jwtCovers.NewJwtManagerCover(repositories.JwtManager{})
	userService := services2.NewUserService(repositories.NewUserRepository(coonPool), producer, jwtManagerCover, tokenService)
	designerProfileService := services.NewDesignerProfileService(repositories2.NewDesignerProfileRepository(coonPool), producer)
	adminService := services3.NewAdminService(repositories3.NewAdminRepository(coonPool))

	serverManager := NewServerManager(
		tokenService,
		userService,
		designerProfileService,
		adminService,
		jwtManagerCover,
		ctx,
		loggerCover,
	)

	serverManager.StartServer()
}
