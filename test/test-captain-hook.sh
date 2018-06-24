#!./test/libs/bats/bin/bats
# Written by Eric Crosson
# 2017-03-12
#
# Test debug method of captain-hook.sh.

load 'libs/bats-support/load'
load 'libs/bats-assert/load'

setup() {
    . "${BATS_TEST_DIRNAME}/../captain-hook.sh"
    TMPDIR=$(mktemp -d bats.XXXX)
}

teardown() {
    [ -d "${TMPDIR}" ] && rm -rf "${TMPDIR}"
}

@test "debug should print nothing when debug mode is not enabled" {
    run debug 'A light switch is also a dark switch.'
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '' ]
}

@test "debug mode should print all args when debug mode is enabled" {
    _debug=1
    debug "Nope"
    truth='When you clean a vacuum cleaner, you are the vacuum cleaner.'
    run debug "${truth}"
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = "${truth}" ]
}

generateRandomValidHookPipeline() {
    prefix="$1"; shift
    hooks="${prefix}/hooks"
    mkdir -p "${hooks}"
    generateRandomValidHookScripts "${hooks}"
}

generateRandomValidHookScripts() {
    outputDir="${1:-${TMPDIR}}"
    touch "${outputDir}/$(( RANDOM % 10))$(( RANDOM % 10))$(( RANDOM % 10))-generated-test.sh"
}

generateRandomInvalidHookPipeline() {
    prefix="$1"; shift
    hooks="${prefix}/hooks"
    mkdir -p "${hooks}"
    generateRandomInvalidHookScripts "${hooks}"
}

generateRandomInvalidHookScripts() {
    outputDir="${1:-${TMPDIR}}"
    # random number between 5 and 10
    numScripts=$(( ( RANDOM % 5 ) + 5 ))
    for i in $(seq 1 "${numScripts}"); do
        touch "${outputDir}/NaN-super-test-should-be-ignored.sh"
        touch "${outputDir}/NaN-123-test-should-be-ignored.sh"
        touch "${outputDir}/NaN-123+-test-should-be-ignored.sh"
    done
}

@test "_getPipeline should return files in numerical order" {
    generateRandomValidHookPipeline "${TMPDIR}"
    run _getPipeline "${TMPDIR}"
    assert [ "${status}" -eq 0 ]
    assert [ "$(echo ${output} | grep -o '[0-9][0-9][0-9]')" = "$(echo ${output} | grep -o '[0-9][0-9][0-9]' | sort -u)" ]
}

@test "_getPipeline should ignore files not beginning with a three-digit number" {
    generateRandomInvalidHookPipeline "${TMPDIR}"
    run _getPipeline "${TMPDIR}"
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '' ]
}

@test "_ensuredScript should return 0 when script is an ensure script" {
    dut="123+-ensured-for-sure.sh"
    run _ensuredScript "${dut}"
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '' ]
}

@test "_ensuredScript should return 1 when script is not an ensure script" {
    dut="123-NOT-ensured.sh"
    run _ensuredScript "${dut}"
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '' ]
}

@test "_updateRc should return $2 when $1 is zero" {
    run _updateRc 0 5
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" -eq 5 ]
}

@test "_updateRc should return $1 when $1 is non-zero" {
    run _updateRc 1 7
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" -eq 1 ]
}

@test "_shortCircuit should return 0 when $1 is equal to zero" {
    run _shortCircuit 0
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '' ]
}

@test "_shortCircuit should return 0 when $1 is less than zero" {
    run _shortCircuit -7
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '' ]
}

@test "_shortCircuit should return 1 when $1 is greater than zero" {
    run _shortCircuit 7
    assert [ "${status}" -eq 1 ]
    assert [ "${output}" = '' ]
}

@test "_determineExitCode should return 0 when $1 is zero" {
    run _determineExitCode 0
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '0' ]
}

@test "_determineExitCode should return 0 when $1 is less than zero" {
    run _determineExitCode -7
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '0' ]
}

@test "_determineExitCode should return $1 when $1 is greater than zero" {
    run _determineExitCode 7
    assert [ "${status}" -eq 0 ]
    assert [ "${output}" = '7' ]
}
