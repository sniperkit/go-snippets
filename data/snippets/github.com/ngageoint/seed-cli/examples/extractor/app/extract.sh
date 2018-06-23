#! /bin/sh
unzip $1 $2 $3
exitCode=$?; if [[ $exitCode != 0 ]]; then exit ${exitCode}; fi
cp seed.outputs.json $3
cat seed.outputs.json
cp seed.png.metadata.json $3
cat seed.png.metadata.json
ls -lR /the
echo ${HELLO}
ls -lR $*
