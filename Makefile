include .env
export

init-system:
	docker-compose up --build

start-system:
	docker-compose up

migrate-up-auth:
	migrate -path "/auth-service/migrations" -database ${DB_CONN_PATH_AUTH} up

migrate-down-auth:
	migrate -path "/auth-service/migrations" -database ${DB_CONN_PATH_AUTH} down

migrate-up-book-cover:
	migrate -path "/book-cover-service/migrations" -database ${DB_CONN_PATH_BOOK_COVER} up

migrate-down-book-cover:
	migrate -path "/book-cover-service/migrations" -database ${DB_CONN_PATH_BOOK_COVER} down