
all:
	GOOS=linux GOARCH=amd64 go build direct_topfunc_call.go

assembly:
	#GOOS=linux GOARCH=amd64 go tool compile -S direct_topfunc_call.go > direct_topfunc_call.as
	go tool compile -S direct_topfunc_call.go > direct_topfunc_call.as

clean:
	rm -f direct_topfunc_call
	rm -f direct_topfunc_call.as
	rm -f direct_topfunc_call.o
