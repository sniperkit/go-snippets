goflags := -v -i -ldflags "-s"
testing :=  -ldflags "-X gitlab.com/poloo/config.Testing=true"

n := qa

local:
	go build $(goflags) -o $(n) $(testing)
	killall $(n) ; ./$(n)&

cpuprof:
	go test -run=XXX -bench=. -benchtime=25s -cpuprofile=cpu.out $(testing)
	go tool pprof $(n).test cpu.out

memprof:
	go test -run=XXX -bench=. -benchtime=25s -memprofile=mem.out $(testing)
	go tool pprof $(n).test mem.out


bench:
	go test -run=XXX -bench=. -benchtime=50s $(testing)


lint:
	gometalinter --disable golint --disable gotype -j 1 ./...

test:
	go test $(testing)


host := 217.182.70.167

upload:
	GOOS=linux GOARCH=amd64 go build $(goflags) -o $(n)_amd64
	strip $(n)_amd64
	ssh -f -n root@$(host) 'mkdir -p /srv/www/poloo.pw/qa/'
	ssh -f -n root@$(host) 'cd /srv/www/poloo.pw/qa/; mv -f $(n) $(n)_old'
	scp $(n)_amd64 root@$(host):/srv/www/poloo.pw/qa/$(n)
	ssh -f -n root@$(host) 'cd /srv/www/poloo.pw/qa/ && (killall $(n) || nohup ./$(n) > output.log 2> errors.log < /dev/null &)'



clean:
	rm -f $(n) $(n)_old $(n)_amd64 cpu.out mem.out profile*.svg
	go clean
