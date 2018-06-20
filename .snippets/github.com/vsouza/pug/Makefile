SERVICE=pug
VERSION=0.0.1

deps:
	@godep restore

build:deps
	@go build -o $(SERVICE)

run-prod:build
	@./$(SERVICE) -env=prod

run-dev:build
	@./$(SERVICE) -env=development

run:build
	@./$(SERVICE)

clean:
	@rm $(SERVICE)
