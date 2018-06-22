#!/bin/sh
set -e
for rel in $(for ver in $(find versions/0.[0-9] | sort); do eval "find versions -name '$(basename $ver).[0-9]*' | sort"; echo $ver; done;); do
    docker push dweomer/vault:$(basename $rel)
done
for rel in $(for ver in $(find versions/0.[0-9][0-9] | sort); do eval "find versions -name '$(basename $ver).[0-9]*' | sort"; echo $ver; done; echo latest); do
    docker push dweomer/vault:$(basename $rel)
done
