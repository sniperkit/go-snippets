#!/usr/bin/env bash

EXTERNAL_DEPS=$(cat <<- END
    golang.org/x/tools/cmd/goimports
    golang.org/x/tools/cmd/cover
    golang.org/x/tour/tree
    honnef.co/go/tools/cmd/megacheck
END
)

for dep in ${EXTERNAL_DEPS}; do
  echo "Installing/Updating ${dep}"
  go get -u "${dep}"
done
