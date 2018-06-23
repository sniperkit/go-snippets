all:
	go fmt src/*.go
	go vet src/*.go
	go build -o gfs src/*.go
install:
	cp gfs /usr/local/bin/
