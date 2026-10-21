-include .env
export

export PROJECT_ROOT=$(shell pwd)




database-start:
	@docker compose up -d data-base

database-down:
	@docker compose down data-base

port-forwarder-postgres-start:
	@docker compose up -d forwarder-port

port-forwarder-postgres-stop:
	@docker compose down forwarder-port

port-forwarder-redis-start:
	@docker compose up -d forwarder-port-redis

port-forwarder-redis-stop:
	@docker compose down forwarder-port-redis

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
		if curl -fsSL http://localhost:8080/health; then \
			echo "App is ready!"; \
			go test ./...; \
			exit 0; \
		fi; \
		echo "App is not ready"; \
		sleep 2; \
  	done && \
  	echo "App failed check" && \
  	exit 1

test-app-server-run:
	@for i in 1 2 3 4 5; do \
			if curl -fsSL "$(IP_SERVER_TEST)/health" > /dev/null; then \
				echo "App is ready!"; \
				go test ./...; \
				exit 0; \
			fi; \
			echo "App is not ready"; \
			sleep 2; \
		done; \
		echo "App failed check"; \
		exit 1

test-domain:
	@go test ./internal/core/domain/ -cover

swagger-generate:
	@docker compose run --rm swagger \
	init \
	-g cmd/trackerapp/main.go \
	-o docs \
	--parseInternal \
	--parseDependency