include .env
export

init-system:
	docker-compose up --build

start-system:
	docker-compose up

migrate-up-auth:
	docker-compose run --rm migrate-auth -path /migrations -database ${DB_CONN_PATH_AUTH} up

migrate-down-auth:
	docker-compose run --rm migrate-auth -path /migrations -database ${DB_CONN_PATH_AUTH} down

migrate-up-book-cover:
	docker-compose run --rm migrate-book-cover -path /migrations -database ${DB_CONN_PATH_BOOK_COVER} up

migrate-down-book-cover:
	docker-compose run --rm migrate-book-cover -path /migrations -database ${DB_CONN_PATH_BOOK_COVER} down

start-test-integration:
	docker-compose up -d postgres-test
	powershell -NoProfile -Command "while ((docker inspect -f '{{.State.Health.Status}}' $$(docker compose ps -q postgres-test)) -ne 'healthy') { Start-Sleep -Seconds 1 }"
	docker-compose run --rm migrate-auth -path /migrations -database ${DB_CONN_PATH_TEST_INT} up
	cd auth-service && go test -v ./...
	docker-compose stop postgres-test
	docker-compose rm -f postgres-test
