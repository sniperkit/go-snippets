PKG=github.com/gophergala/tron

ec2:
	git archive --output=ec2.zip HEAD

local:
	rm -rf ${GOPATH}/pkg/darwin_amd64/${PKG}
	go get -tags local ${PKG}/bin/server
	ASSETS_PATH=frontend/src ${GOPATH}/bin/server -logtostderr=true -stderrthreshold=INFO

clean:
	rm -rf ec2.zip
