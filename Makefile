# Compose
up:
	docker compose up -d --build

down:
	docker compose down

build-broker:
	docker build -d broker-service .

# Frontend
run:
	go run /front-end/main.go