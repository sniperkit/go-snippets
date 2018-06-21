protofiles:
	# Go
	protoc -I/usr/local/include -I. \
  		-I ${GOPATH}/src \
 		-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
 		--go_out=google/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
  		proto/*.proto
	# Gateway
	protoc -I/usr/local/include -I. \
   		-I ${GOPATH}/src \
   		-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
   		--grpc-gateway_out=logtostderr=true:. \
   		proto/*.proto
	# Swagger
	protoc -I/usr/local/include -I. \
  		-I ${GOPATH}/src \
  		-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  		--swagger_out=logtostderr=true:. \
  		proto/*.proto
run-api:
	GRPC_PORT=50051 HTTP_PORT=8080 ELASTICSEARCH_URL=localhost:9200 go run api/cmd/api/main.go
