package main

import (
	logs "book-cover-service/core/logger"
	"book-cover-service/core/middleware"
	"book-cover-service/internal/auth/repositories"
	services2 "book-cover-service/internal/books/services"
	"book-cover-service/internal/books/transport"
	services3 "book-cover-service/internal/covers/services"
	transport2 "book-cover-service/internal/covers/transport"
	services4 "book-cover-service/internal/reactions/services"
	transport3 "book-cover-service/internal/reactions/transport"
	"book-cover-service/internal/snapshots/services"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerManager struct {
	coverHandlers     transport2.HTTPCoverHandlers
	coverLikeHandlers transport3.HTTPCoverLikeHandlers
	favoritesHandlers transport3.HTTPFavoritesHandlers
	bookHandlers      transport.HTTPBookHandlers
	middleware        middleware.AuthMiddleware
}

func NewServerManager(coverService services3.CoverService,
	coverLikeService services4.CoverLikeService,
	favoritesService services4.FavoritesService,
	designerProfileSnapshotService services.DesignerProfileSnapshotService,
	bookService services2.BookService,
	JWTManager repositories.JWTManager,
	ctx context.Context,
	logger logs.Logger) ServerManager {
	return ServerManager{
		coverHandlers:     transport2.NewCoverHandlers(coverService, designerProfileSnapshotService, ctx, logger),
		coverLikeHandlers: transport3.NewCoverLikeHandlers(coverLikeService, ctx, logger),
		favoritesHandlers: transport3.NewHTTPFavoritesHandlers(favoritesService, ctx, logger),
		bookHandlers:      transport.NewHTTPBookHandlers(bookService, ctx, logger),
		middleware:        middleware.NewAuthMiddleware(JWTManager),
	}
}

func (s *ServerManager) StartServer() {
	router := mux.NewRouter()

	router.
		Path("/books").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.bookHandlers.HandleAddBook)),
		)

	router.
		Path("/books").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.bookHandlers.HandleUpdateBook)),
		)

	router.
		Path("/books/{book_id}").
		Methods("DELETE").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.bookHandlers.HandleDeleteBook)),
		)

	router.
		Path("/designers/me/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.coverHandlers.HandleGetMyCoversAsDesigner)),
		)

	router.
		Path("/designers/{user_id}/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.Auth(
				http.HandlerFunc(s.coverHandlers.HandleGetCoversByDesignerUserID),
			),
		)

	router.
		Path("/covers/{cover_id}").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.Auth(
				http.HandlerFunc(s.coverHandlers.HandleUpdateCover),
			),
		)

	router.
		Path("/covers/{cover_id}").
		Methods("GET").
		HandlerFunc(s.coverHandlers.HandleGetCoverByID)

	router.
		Path("/covers").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.coverHandlers.HandleAddCover)),
		)

	router.
		Path("/feeds/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(s.coverHandlers.HandleGetFeedCovers)

	router.
		Path("/books/{book_id}/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(s.coverHandlers.HandleGetCoversByBook)

	router.
		Path("/covers/{cover_id}/like").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.coverLikeHandlers.HandleSetLike)),
		)

	router.
		Path("/covers/{cover_id}/favorite").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.favoritesHandlers.HandleAddCoverToFavorites)),
		)

	router.
		Path("/covers/me/favorites").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.favoritesHandlers.HandleGetMyFavoriteCovers)),
		)

	err := http.ListenAndServe(":9011", router)
	if err != nil {
		println("Ошибка при запуске сервера")
		return
	}
}
