.PHONY: dev build up down clean

dev:
	docker compose up --build

build:
	docker compose build

up:
	docker compose up

down:
	docker compose down
	sudo rm -rf tmp

clean: down
	docker compose rm -f
	docker rmi frfr-app
	sudo rm -rf tmp
