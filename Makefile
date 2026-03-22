build:
	@go build -o bin/app ./cmd

run:build
	@./bin/app

push:
	@git init
	@git add .
	@git commit -s -m"${msg}"
	@echo "pushing all files to git repository..."
	@git push

swag:
	@swag init -g cmd/main.go

create-migration:
	@migrate create -ext=sql -dir=internal/database/migrations -seq init

migrate-up-postgres:
	@set -a; . ./.env.development; set +a; \
	migrate -path=internal/database/migrations \
	-database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DBNAME?sslmode=$$DB_SSL_MODE" \
	-verbose up


migrate-down-postgres:
	@set -a; . ./.env.development; set +a; \
	migrate -path=internal/database/migrations \
	-database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DBNAME?sslmode=$$DB_SSL_MODE" \
	-verbose down

migrate-reset:
	@set -a; . ./.env.development; set +a; \
	migrate -path internal/database/migrations \
	-database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DBNAME?sslmode=$$DB_SSL_MODE" \
	force $(VERSION)

migrate-drop:
	@set -a; . ./.env.development; set +a; \
	migrate -path internal/database/migrations \
	-database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DBNAME?sslmode=$$DB_SSL_MODE" \
	drop -f