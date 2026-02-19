build:
	@go build -o bin/app .

run:build
	@./bin/app

push:
	@git init
	@git add .
	@git commit -s -m"${msg}"
	@echo "pushing all files to git repository..."
	@git push

create-migration:
	@migrate create -ext=sql -dir=internal/database/migrations -seq init

migrate_up_postgres:
	@set -a; . ./.env.development; set +a; \
	migrate -path=internal/database/migrations \
	-database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DBNAME?sslmode=$$DB_SSL_MODE" \
	-verbose up


migrate_down_postgres:
	@set -a; . ./.env.development; set +a; \
	migrate -path=internal/database/migrations \
	-database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DBNAME?sslmode=$$DB_SSL_MODE" \
	-verbose down
