#!/bin/sh
set -eo pipefail

# constants
readonly ARGS="$@"

# includes
readonly UTILS_FILE='/usr/local/bin/entrypoint-utils.sh'
if [[ -f "${UTILS_FILE}" ]]; then
    source "${UTILS_FILE}"
else
    echo -e "${UTILS_FILE} is missing!" >&2
    exit 1
fi

# funcs
parse_cmd() {
    # debug
    #info "${FUNCNAME}" "$*"

    local _CMD="$1"

    [[ -n ${_CMD} && -x ${_CMD} ]] \
        || fatal 'Please specify a program to run.'

    shift
    local _ARGS="$@"

    # run as specified user
    exec su-exec \
        "${RUN_AS_USER}" \
        "${_CMD}" \
        ${_ARGS}
}

main() {
    # env check
    [[ -n ${RUN_AS_USER} ]] \
        || fatal "This script is only compatible with project nutshells images,\ndon't try to run it out straight."

    parse_cmd ${ARGS}
}

trap cleanup EXIT
main
