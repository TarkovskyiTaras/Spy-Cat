.PHONY: all up down run build migrate

up:
	docker-compose -f deploy/local/docker-compose.yml up -d db

down:
	docker-compose -f deploy/local/docker-compose.yml down

run:
	go run ./cmd/app

all start: up run
