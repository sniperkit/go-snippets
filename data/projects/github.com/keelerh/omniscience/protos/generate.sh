#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
docker_dir="${__dir}/build"


docker build -t omniscience-proto-builder "${docker_dir}"
docker run -v "${__dir}:/protos" -w "/protos" omniscience-proto-builder /go/generate_protos.sh
