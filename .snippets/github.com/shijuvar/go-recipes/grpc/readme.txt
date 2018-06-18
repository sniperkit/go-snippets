
To compile with protoc, run the following command from root directory (grpc directory):

protoc -I customer/ customer/customer.proto --go_out=plugins=grpc:customer

