#!/usr/bin/env bash

set -e

# Ensure script directory is CWD
pushd "${0%/*}" > /dev/null


UNAME=$(uname -s)

echo Building example images.................................................................
SEED=""
case "${UNAME}" in
    Linux*)           SEED=output/seed-linux-amd64; SUDO=sudo;;
    Darwin*)          SEED=output/seed-darwin-amd64;;
    CYGWIN*)          SEED=output/seed-windows-amd64;;
    MINGW64_NT-10.0*) SEED=output/seed-windows-amd64;;
    *)                SEED="UNKNOWN:${UNAME}"
esac
${SUDO} ${SEED} build -d examples/addition-job/
${SUDO} ${SEED} build -d examples/extractor/

echo Finished building example images.................................................................

echo Running example images...........................................................................
${SUDO} ${SEED} run -in addition-job-0.0.1-seed:1.0.0 -i INPUT_FILE=examples/addition-job/inputs.txt -rm -m MOUNT_BIN=testdata/complete/ -m MOUNT_TMP=testdata/ -e SETTING_ONE=one -e SETTING_TWO=two -o temp
echo ...
echo ...
${SUDO} ${SEED} run -in extractor-0.1.0-seed:0.1.0 -i ZIP=testdata/seed-scale.zip -i MULTIPLE=examples/extractor/seed.manifest.json -i MULTIPLE=examples/extractor/seed.outputs.json -rm -m MOUNTAIN=testdata/complete/ -e HELLO=Hello -o temp
echo Finished running example images..................................................................

popd >/dev/null
