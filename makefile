include .env
export

export PROJECT_ROOT=$(shell pwd)

db-up:
	@docker compose up -d svc-postgres

db-down:
	@docker compose down svc-postgres

db-acl:
	@sudo setfacl -R -m u:$(USER):rwx ./out/pgdata

db-clean:
	@make db-down && sudo rm -rf ./out/pgdata

db-migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make db-migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm svc-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

db-migrate-up:
	@make db-migrate-action action=up

db-migrate-down:
	@make db-migrate-action action=down

db-migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make db-migrate-action action=1"; \
		exit 1; \
	fi; \
	docker compose run --rm svc-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@svc-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"
