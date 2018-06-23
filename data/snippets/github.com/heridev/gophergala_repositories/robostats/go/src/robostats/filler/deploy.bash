#/usr/bin/env bash
(cd ../../../ && source ./env.bash)
go build -o filler
scp filler www-data@dev.robostats.io:
ssh www-data@dev.robostats.io './filler'
rm filler
