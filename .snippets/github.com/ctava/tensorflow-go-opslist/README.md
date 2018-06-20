# Installation
`git clone https://github.com/ctava/tensorflow-go-opslist.git`

Run the following commands:
```
cd tensorflow-go-opslist
docker build -t tensorflow-go-opslist .
docker run -p 8888:8888 -d tensorflow-go-opslist
docker ps
```
Example output:

| CONTAINER ID  | IMAGE                 | 
| ------------- | --------------------- |
| e726f3ee010c  | tensorflow-go-opslist | 
 

take the CONTAINER ID and add it to the following command:
```
docker exec -it <CONTAINER ID> bash
```
and your in. You now have `tensorflow` + `golang` available.

# Confirm Golang and Tensorflow Installation

Run the following commands:
```
go version
cd src/github.com/ctava/tensorflow-go-version
go run main.go
```

# Now for the fun stuff - get a list of operations
```
cd ../tensorflow-go-opslist
go run main.go
```
You should see a list of all of the ops available in tensorflow

# Additional resources
[gopherdata](https://github.com/gopherdata/resources)
