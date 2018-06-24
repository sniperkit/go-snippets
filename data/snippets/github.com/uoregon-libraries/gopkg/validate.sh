#!/usr/bin/env bash
#
# validate.sh runs various code analysis tools to reduce the likelihood of
# unnoticed errors
set -eu

err() {
  echo $@ >&2
}

IFS=''
unformatted=$(find ./ -name "*.go" | xargs gofmt -l -s)
linter=$(golint ./...)
vet=$(go vet -printfuncs Debugf,Infof,Warnf,Errorf,Criticalf,Fatalf ./...  2>&1 || true)

result=0

if [[ $unformatted != "" ]]; then
  err "gofmt reports issues:"
  err "---------------------"
  err $unformatted
  result=1
fi

if [[ $linter != "" ]]; then
  if [[ $result == 1 ]]; then
    err
  fi
  err "linter reports issues:"
  err "----------------------"
  err $linter
  result=1
fi

if [[ $vet != "" ]]; then
  if [[ $result == 1 ]]; then
    err
  fi
  err "go vet reports issues:"
  err "----------------------"
  err $vet
  result=1
fi

if [[ $result == 1 ]]; then
  err
fi
exit $result
