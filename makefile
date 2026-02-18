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