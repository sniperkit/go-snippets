#!/usr/bin/env bash
# Written by Eric Crosson
# 2016-10-28
#
# This script provides a framework for git hooks.
#
# Goals:
# - provide an isolated space for each hook
# - provide a scalable framework for development to expand into
# - DRY
# - handle edge cases in git flows (mainly due to stashing)
#
# Note: functions beginning with an underscore are private.

# set -x

#------------------------------------------------------------------------------
# Define constants
#------------------------------------------------------------------------------
# Set this flag to a 1 to get debugging output
#
# Note: may not work great with the in-place test status updates.
#       Consider logging with `tee` to a file for inspection.
declare _debug=0

# Use this command instead of `echo` to only print information when
# debug-mode is enabled.
debug() {
    if [[ "${_debug}" -gt 0 ]]; then
        >&2 echo "$@"
    fi
}

# Return the sorted list of files comprising a git-hook pipeline.
_getPipeline() {
    searchDir=$1; shift
    find "${searchDir}" \
         -regex ".*/[0-9][0-9][0-9]+?.*" \
        | sort
}

# Test if the supplied script is classified as a "ensured" script.
_ensuredScript() {
    [ -z "$(echo "${script}" | sed -ne '/[0-9]\{3\}+/p')" ]
}

# Given the existing exit-code ($1) and the exit-code of the last
# script that ran ($2), return the new exit-code for this pipeline.
_updateRc() {
    if [[ $1 -eq 0 ]]; then
        echo "$2"
    else
        echo "$1"
    fi
}

# Test if this pipeline is in short-circuit mode, given the current
# pipeline's exit-code.
_shortCircuit() {
    local -r rc=$1; shift
    [ "${rc}" -le 0 ]
}

# Determine the exit code for the entire pipeline.
#
# Note: this function should be called at the end of a pipeline's
# execution.
_determineExitCode() {
    local -r rc=$1; shift
    # negative exit codes indicate a short-circuit without failure
    if _shortCircuit "${rc}"; then
        echo 0
    else
        echo "${rc}"
    fi
}

# This method is called by captain-hook.  It will run the user's
# defined git-hook pipeline when a client-side hook is invoked by git.
runHook() {
    local -r executing_hook=$1; shift
    local rc=0

    local -r _tmpfile=$(mktemp captain-hook.XXXX)
    _getPipeline "${executing_hook}.d" > "${_tmpfile}"
    debug "PWD is $PWD"
    debug "Pipeline is $(cat "${_tmpfile}")"
    while IFS= read -r script; do
        # if part of the pipeline has failed, skip forward to the
        # "ensured" scripts
        if [[ "${rc}" -ne 0 ]]; then
            ! _ensuredScript "${script}" && continue
        fi

        debug "Executing script: '${script}'..."
        "${script}" "$@"
        rc=$(_updateRc "${rc}" $?)
        debug "Executing script: '${script}'...done"
    done < "${_tmpfile}"
    rm -f "${_tmpfile}"

    exit "$(_determineExitCode "${rc}")"
}
export runHook
