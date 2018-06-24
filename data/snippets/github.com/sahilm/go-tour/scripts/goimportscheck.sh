#!/usr/bin/env bash

goimports_files="$(goimports -l ./)"
if [[ -n ${goimports_files} ]]; then
    echo 'goimports needs running on the following files:'
    echo "${goimports_files}"
    echo "You can use the command: \`make goimports\` to reformat code."
    exit 1
fi
