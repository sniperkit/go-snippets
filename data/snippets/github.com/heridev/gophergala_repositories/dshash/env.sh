pwd=`pwd`
go=.goroot

export PATH=$pwd/bin:$pwd/$go:$PATH:$pwd/$go/goroot/bin
export GOROOT=$pwd/$go/goroot
export GOPATH=$pwd/.vendor:$pwd
