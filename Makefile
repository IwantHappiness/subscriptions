include .env
export

export PROJECT_ROOT := $(shell pwd)

env-build:
	@docker-compose build

env-up:
	@docker-compose up -d app

env-down:
	@docker-compose down postgres

migrate-create:
	$(if $(seq),,$(error seq is not set. Usage: make migrate-create seq=name))
	@docker-compose run --rm migrate \
		create \
		-ext sql \
		-dir ./migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down


migrate-action:
	@docker-compose run --rm migrate \
		-path ./migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(action)"
