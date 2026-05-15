include .env
export

export PROJECT_ROOT=$(shell pwd)


env-up:
	@docker compose up -d data-base

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down data-base forwarder-port && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

logs-clean:
	@read -p "Очистить все логи? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "логи очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Dont have seq"; \
		exit 1; \
	fi; \
		docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Dont have action"; \
		exit 1; \
	fi; \
	docker compose run --rm migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@data-base:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

clean_migrate:
	@docker compose run --rm migrate \
		-path /migrations \
     	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@data-base:5432/${POSTGRES_DB}?sslmode=disable force 0

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

database-start:
	@docker compose up -d data-base

database-down:
	@docker compose down data-base

port-forwarder-start:
	@docker compose up -d forwarder-port

port-forwarder-stop:
	@docker compose down forwarder-port

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/trackerapp/main.go

tracker-deploy-run:
	@docker compose up -d --build tracker-app

tracker-deploy-stop:
	@docker compose down tracker-app

tracker-deploy-check:
	@docker compose up -d --build tracker-app
	@for i in 1 2 3 4 5; do \
		if curl -fsSL http://localhost:8080/; then \
			exit 0; \
		fi; \
		echo "App is not ready"; \
		sleep 2; \
  	done && \
  	echo "App failed check" && \
  	exit 1

test-app-server-run:
	@for i in 1 2 3 4 5; do \
		if curl -fsSL ${IP_SERVER_TEST}; then \
			exit 0; \
		fi; \
		echo "App is not ready"; \

		sleep 2; \
  	done && \
  	echo "App failed check" && \
  	exit 1

swagger-generate:
	@docker compose run --rm swagger \
	init \
	-g cmd/trackerapp/main.go \
	-o docs \
	--parseInternal \
	--parseDependency