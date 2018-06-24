#!/usr/bin/env sh
# set -e

dep ensure # This is required by docker-compose and hot-reloading
dep ensure -update
fresh -c scripts/runner.conf
