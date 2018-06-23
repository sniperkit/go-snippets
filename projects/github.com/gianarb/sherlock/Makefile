sherlockd/build:
	go build -o ./bin/sherlockd backend/cmd/sherlockd/main.go
	GOOS=darwin GOARCH=amd64 go build -i -o ./bin/sherlockd_darwin backend/cmd/sherlockd/main.go
