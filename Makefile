include .env

install_deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

update_sql:
	sqlc generate

down: #migrate down
	cd sql/schema && goose postgres ${DATABASE_URL} down

up: #migrate up
	cd sql/schema && goose postgres ${DATABASE_URL} up