build:
	go build -o ./bin/plugserver

run: build
	./bin/plugserver
