default: all

dep:
	@cd cmd/ && go get -d -v

all: dep
	go build cmd/spider.go
	go build cmd/server.go

tools:
	go build cmd/test.go
	go build cmd/print.go
	go build cmd/search.go

tests:
	go test ./database/ ./models/ ./retrieval/ ./stopword/ -cover

tests_report:
	go test ./database/ ./models/ ./retrieval/ ./stopword/ -coverprofile=c.out
	go tool cover -html=c.out

clean:
	rm -f spider test server search phase1.zip print
