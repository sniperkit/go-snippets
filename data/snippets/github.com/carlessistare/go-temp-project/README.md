## Install
* Install Go on your computer (https://golang.org/doc/install)
* Add Go binary in your local path
```bash
PATH=$PATH:$GOPATH/bin
```

* Install Godep and download dependencies
```bash
go get github.com/tools/godep
godep restore
```

## Before commit on github
If you add new packages you have to update the Godeps.json by doing this:
```bash 
godep save ./...
```
_Info: It's important to put "./..." after the "save" to include test dependencies.
 (The parameter "-t" exist in Godep to include test files but it does not work)_

### Launch tests
```bash
go test ./test/... 
``` 
