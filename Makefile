docker:
	docker compose up --build

client:
	go run cmd/client/main.go

server:
	go run cmd/server/main.go

test:
	go test ./... -count 1
