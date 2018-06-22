.PHONY:	build clean fmt test vet

APP_NAME = Project-V
EXEC_NAME = ./$(APP_NAME)

$(EXEC_NAME):
	go build

build: $(EXEC_NAME) 

fmt:
	go fmt .

test:
	go test .

dep:
	go get

clean:
	go clean
