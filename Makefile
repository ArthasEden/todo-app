include .env
export

export PROJECT_ROOT=${shell pwd}

# db

# Запустить контейнер svc-db
db-up:
	@docker compose up -d svc-db && make out-permission

# Остановить и удалить контейнер svc-db
db-down:
	@docker compose down svc-db 

# Остановить и удалить контейнер svc-db, svc-db-port для последующей очистки папки `volumes`
db-clean:
	@docker compose down svc-db && docker compose down svc-db-port && sudo rm -rf ./out/pgdata

# db-port

# Запустить контейнер svc-db-port
db-port-up:
	@docker compose up -d svc-db-port

# Остановить и удалить контейнер svc-db-port
db-port-down:
	@docker compose down svc-db-port


# db-migrate

# Запустить контейнер с самоудалением svc-db-migrate для последующего создания файлов миграции
# Здесь `run --rm` как типичный паттерн для одноразовых контейнеров.
# Пример: make db-migrate-create seq=init
db-migrate-create:
	@docker compose run --rm svc-db-migrate create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)" \
		&& make db-migrate-permission

# Выдать права текущему пользователю на просмотр и изменение папки `migrate`
db-migrate-permission:
	@sudo chown -R $$USER:$$USER migrations

# Внутренняя команда для выполнения команды migrate
db-migrate-action:
	@docker compose run --rm svc-db-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@svc-db:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)" $(args)

# Запустить контейнер с самоудалением svc-db-migrate для последующего применениия миграции
db-migrate-up:
	@make db-migrate-action action=up

# Запустить контейнер с самоудалением svc-db-migrate для последующего отката миграции
db-migrate-down:
	@make db-migrate-action action=down

# Принудительно устанавливает номер версии миграции и сбрасывает статус миграции `dirty`.
# Пример: db-migrate-force version=1
db-migrate-force:
	@make db-migrate-action action=force args="$(version)"


# other

# Выдать права текущему пользователю на просмотр и изменение содержимого папки `out`
out-permission:
	@sudo setfacl -R -m u:$(USER):rwx ./out
# Запуск приложения
run:
	@go mod tidy && go run cmd/main.go