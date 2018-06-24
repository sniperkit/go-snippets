# COMP 4321 Search Engine [![Go Report Card](https://goreportcard.com/badge/github.com/rsmohamad/comp4321)](https://goreportcard.com/report/github.com/rsmohamad/comp4321)

## Getting started
    
 - Install Go:
   ```
   wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
   sudo tar -C /usr/local/ -xzf go1.10.linux-amd64.tar.gz
   echo "export PATH=\$PATH:/usr/local/go/bin" | sudo tee -a /etc/profile
   source /etc/profile
   ```
   
- Make sure the `$GOPATH` environment variable is set to `~/go/`. Check with `go env`.

- Download the repository
    ```
    go get github.com/rsmohamad/comp4321
    ```


## Building

- Inside the project directory, type `make`
- `./spider [-start=<starting page>] [-pages=<number of pages>] [-a] ` to run the spider
- `./server` to launch the webserver

