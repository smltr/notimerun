dev:
	docker compose up --build

devlog:
	docker compose up --build > findservers.log 2>&1

build-dev:
	docker compose build

up-dev:
	docker compose up

down:
	docker compose down
	sudo rm -rf tmp

clean: down
	docker compose rm -f
	docker rmi findservers-app
	sudo rm -rf tmp

# Deploy to fly.io
deploy:
	fly deploy
