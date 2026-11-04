api = apps/api
extension = apps/extension

.PHONY: up down dev migrate test fmt vet build check

up:
	docker compose up --build

down:
	docker compose down

dev:
	cd $(api) && go run ./cmd/api

migrate:
	cd $(api) && go run ./cmd/migrate

test:
	cd $(api) && go test ./...

fmt:
	cd $(api) && gofmt -w .

vet:
	cd $(api) && go vet ./...

build:
	cd $(api) && go build ./...
	cd $(extension) && npm run build

check:
	@cd $(api) && test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	cd $(api) && go vet ./...
	cd $(api) && go test ./...
	cd $(extension) && npm run typecheck
