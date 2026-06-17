ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: run migrate-up migrate-down swag

run:
	go run cmd/app/main.go

migrate-up:
	migrate -path migrations -database "${PG_URL}" up

migrate-down:
	migrate -path migrations -database "${PG_URL}" down

migrate-force:
	migrate -path migrations -database "${PG_URL}" force $(version)

swag:
	swag init -g cmd/app/main.go
