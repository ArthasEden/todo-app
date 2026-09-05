include .env
export

export PROJECT_ROOT=$(shell pwd)

# Запустить контейнер PostgreSQL
db-up:
	@docker compose up -d svc-db

# Остановить и удалить контейнер PostgreSQL
db-down:
	@docker compose down svc-postgres

# Выдать права текущему пользователю на просмотр и изменение содержимого volume
db-acl:
	@sudo setfacl -R -m u:$(USER):rwx ./out/pgdata

# Остановить и удалить контейнер PostgreSQL и очистить volumes
db-clean:
	@make db-down && sudo rm -rf ./out/pgdata

# Создать новую SQL-миграцию
# Пример: make db-migrate-create seq=init
db-migrate-create:
	@docker compose run --rm svc-db-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

# Применить SQL-миграцию
db-migrate-up:
	@make db-migrate-action action=up

# Откатить SQL-миграцию
db-migrate-down:
	@make db-migrate-action action=down

# Внутренняя вспомогательная команда для применения/откатки SQL-миграции
# Пример: make db-migrate-action action=up
db-migrate-action:
	@docker compose run --rm svc-db-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@svc-db:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

# Выдать права текущему пользователю на просмотр и изменение содержимого migrations
db-migrate-chown:
	@sudo chown -R $$USER:$$USER migrations

# Принудительно устанавливает номер версии миграции и сбрасывает dirty.
# Пример: make db-migrate-force version=1
db-migrate-force:
	@docker compose run --rm svc-db-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@svc-db:5432/${POSTGRES_DB}?sslmode=disable" \
		force $(version)
# Подключаем порт 5432 для доступа к БД
# Создаём контенер с сервисом Socat
db-port-on:
	@docker compose up -d svc-db-port

# Удаляем контенер с сервисом Socat
db-port-off:
	@docker compose down svc-db-port