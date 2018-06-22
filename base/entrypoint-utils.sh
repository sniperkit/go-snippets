#!/bin/sh
set -eo pipefail

# flags
[[ -n "${UTILS_LOADED}" ]] && return
export UTILS_LOADED='true'

# funcs
throw_exception() {
    local _LEVEL="${1:-INFO}"
    shift
    echo -e "\n[${_LEVEL}]\n$(date +'%Y-%m-%dT%H:%M:%S%z')\n$*" >&2
}

info() {
    throw_exception 'INFO' "$*"
}

warning() {
    throw_exception 'WARNING' "$*"
}

error() {
    throw_exception 'ERROR' "$*"
}

fatal() {
    throw_exception 'FATAL' "$*"
    exit 1
}

get_flag(){
    local _FLG_NAME="${1:-unknown}"
    if [[ "${#_FLG_NAME}" -eq 1 ]] ; then
        echo "-${_FLG_NAME}"
    else
        echo "--${_FLG_NAME}"
    fi
}

cleanup() {
    info 'Entrypoint Script ended.'
}
