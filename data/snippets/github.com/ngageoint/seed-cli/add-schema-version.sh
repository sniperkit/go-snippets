#!/usr/bin/env bash

## Usage: ./add-schema-version.sh VERSION [base_url]

VERSION=$1

if [[ "${VERSION}x" == "x" ]]
then
    echo Missing version parameter!
    echo Usage:
    echo   $0 0.0.0 [base_url]
    exit 1
fi

BASE_URL=https://github.com/ngageoint/seed/releases/download
if [[ "${2}x" != "x" ]]
then
    BASE_URL=$2
    echo Updated base URL to $2
fi

# Ensure script directory is CWD
pushd "${0%/*}"

rm -fr schema/${VERSION}
mkdir -p schema/${VERSION}
cd schema/${VERSION}

wget ${BASE_URL}/${VERSION}/seed.manifest.example.json
wget ${BASE_URL}/${VERSION}/seed.manifest.schema.json
wget ${BASE_URL}/${VERSION}/seed.metadata.schema.json

popd > /dev/null
