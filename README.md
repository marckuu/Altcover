Маркетплейс альтернативных обложек для книг.

Часто многие книги продаются в нескольких не самых красивых вариантах. 
Данная система предлагает пользователям создавать свои обложки для разных книг и публиковать их.


КАК ЗАПУСТИТЬ ПРОЕКТ:
1. В корне проекта выполнить make init-system
2. В папке auth-service выполнить make migrate-up-auth (миграции для auth сервиса)
3. В папке book-cover-service выполнить make migrate-up-book-cover (миграции для book-cover сервиса)
4. Открыть http://localhost:9090/swagger/index.html#/ для документации auth сервиса и http://localhost:9091/swagger/index.html#/ для документации book-cover сервиса

КАК ЗАПУСТИТЬ ИНТЕГРАЦИОННЫЕ ТЕСТЫ в auth сервисе:
1. В папке auth-service выполнить start-test-integration:
