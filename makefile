dev:
	docker compose up --build

build-dev:
	docker compose build

up-dev:
	docker compose up

down:
	docker compose down
	sudo rm -rf tmp

clean: down
	docker compose rm -f
	docker rmi notimerun-app
	sudo rm -rf tmp

# Production commands
build-prod:
	docker build -t notimerun:prod -f Dockerfile.prod .

# Test production build locally
run-prod:
	docker run -p 8080:8080 -e STEAM_API_KEY=${STEAM_API_KEY} notimerun:prod

# Deploy to fly.io
deploy:
	fly deploy
