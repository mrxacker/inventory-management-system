COMPOSE_FILE=docker-compose.yml

ENV_FILE=.env

ifneq (,$(wildcard .env))
	include .env
	export
endif

up:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

down:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down -v