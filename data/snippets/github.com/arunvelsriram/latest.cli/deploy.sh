#! /bin/bash

set -e

curl -sL https://git.io/goreleaser | bash -s -- --config .release.yml --release-notes release-notes/$TRAVIS_TAG.md --debug
