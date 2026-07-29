package main

import (
	"book-cover-service/core/enums"
	logs "book-cover-service/core/logger"
	"book-cover-service/core/middleware"
	"book-cover-service/internal/auth/repositories"
	bookServicesInterfaces "book-cover-service/internal/books/services/interfaces"
	"book-cover-service/internal/books/transport"
	coverServicesInterfaces "book-cover-service/internal/covers/services/interfaces"
	transport2 "book-cover-service/internal/covers/transport"
	reactionServicesInterfaces "book-cover-service/internal/reactions/services/interfaces"
	transport3 "book-cover-service/internal/reactions/transport"
	designerProfileSnapshotInterfaces "book-cover-service/internal/snapshots/services/interfaces"
	"context"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "book-cover-service/docs"

	"github.com/gorilla/mux"
)

type ServerManager struct {
	coverHandlers     transport2.HTTPCoverHandlers
	coverLikeHandlers transport3.HTTPCoverLikeHandlers
	favoritesHandlers transport3.HTTPFavoritesHandlers
	bookHandlers      transport.HTTPBookHandlers
	middleware        middleware.AuthMiddleware
	logger            logs.Logger
}

func NewServerManager(coverService coverServicesInterfaces.CoverService,
	coverLikeService reactionServicesInterfaces.CoverLikeService,
	favoritesService reactionServicesInterfaces.FavoritesService,
	designerProfileSnapshotService designerProfileSnapshotInterfaces.DesignerProfileSnapshotService,
	bookService bookServicesInterfaces.BookService,
	JWTManager repositories.JWTManager,
	ctx context.Context,
	logger logs.Logger) ServerManager {
	return ServerManager{
		coverHandlers:     transport2.NewCoverHandlers(coverService, designerProfileSnapshotService, ctx, logger),
		coverLikeHandlers: transport3.NewCoverLikeHandlers(coverLikeService, ctx, logger),
		favoritesHandlers: transport3.NewHTTPFavoritesHandlers(favoritesService, ctx, logger),
		bookHandlers:      transport.NewHTTPBookHandlers(bookService, ctx, logger),
		middleware:        middleware.NewAuthMiddleware(JWTManager),
		logger:            logger,
	}
}

func (s *ServerManager) StartServer() {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.
		Path("/books").
		Methods("POST").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Admin,
				http.HandlerFunc(s.bookHandlers.HandleAddBook)),
		)

	router.
		Path("/books").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Admin,
				http.HandlerFunc(s.bookHandlers.HandleUpdateBook)),
		)

	router.
		Path("/books/{book_id}").
		Methods("DELETE").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Admin,
				http.HandlerFunc(s.bookHandlers.HandleDeleteBook)),
		)

	router.
		Path("/books/{title}").
		Methods("GET").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.bookHandlers.HandleGetBookByTitle)),
		)

	router.
		Path("/designers/me/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Designer,
				http.HandlerFunc(s.coverHandlers.HandleGetMyCoversAsDesigner)),
		)

	router.
		Path("/designers/{user_id}/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.coverHandlers.HandleGetCoversByDesignerUserID)),
		)

	router.
		Path("/covers/{cover_id}").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Designer,
				http.HandlerFunc(s.coverHandlers.HandleUpdateCover)),
		)

	router.
		Path("/covers/{cover_id}").
		Methods("GET").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.coverHandlers.HandleGetCoverByID)),
		)

	router.
		Path("/covers").
		Methods("POST").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Designer,
				http.HandlerFunc(s.coverHandlers.HandleAddCover)),
		)

	router.
		Path("/feeds/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.coverHandlers.HandleGetFeedCovers)),
		)

	router.
		Path("/books/{book_id}/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.coverHandlers.HandleGetCoversByBook)),
		)

	router.
		Path("/covers/{cover_id}/like").
		Methods("POST").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.coverLikeHandlers.HandleSetLike)),
		)

	router.
		Path("/covers/{cover_id}/favorite").
		Methods("POST").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.favoritesHandlers.HandleAddCoverToFavorites)),
		)

	router.
		Path("/covers/me/favorites").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.favoritesHandlers.HandleGetMyFavoriteCovers)),
		)

	err := http.ListenAndServe(":9011", router)
	if err != nil {
		s.logger.Error(fmt.Errorf("ошибка при запуске сервера, %w", err).Error())
		return
	}
}
