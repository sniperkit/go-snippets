#!/bin/sh
set -eo pipefail

# constants
readonly ARGS="$@"
readonly ENTRYPOINT_CMD='/usr/local/sbin/dnsmasq'

# includes
readonly UTILS_FILE='/usr/local/bin/entrypoint-utils.sh'
if [[ -f "${UTILS_FILE}" ]]; then
    source "${UTILS_FILE}"
else
    echo -e "${UTILS_FILE} is missing!" >&2
    exit 1
fi

# funcs
halt(){
    local _FLG
    _FLG="$(get_flag "$2")"
    local _ERR_MSG="Option '${_FLG}' is managed by the entrypoint script.\n"
    local _ERR_NUM="${1:-0}"

    case "${_ERR_NUM}" in
        1)
            # no forking
            _ERR_MSG="${_ERR_MSG}Pidfile isn't necessary for it'll keep running in foreground.\n"
            ;;
        2)
            # built in
            _ERR_MSG="${_ERR_MSG}Use the '-e', '--env' or '--env-file' options to override its value,\n"
            _ERR_MSG="${_ERR_MSG}referer to README for the relevant environment variables.\n"
            ;;
    esac

    _ERR_MSG="${_ERR_MSG}Just simply remove it and try again."

    fatal "${_ERR_MSG}"
}

run_cmd() {
    local _EXEC="$1"
    shift
    local _ARGS="$@"

    if [[ "${_EXEC}" == 'true' ]]; then
        exec "${ENTRYPOINT_CMD}" \
            -u root \
            ${_ARGS}
    else
        "${ENTRYPOINT_CMD}" \
            -u root \
            ${_ARGS}
    fi
}

check_opts() {
    local _ARGS="$@"

    # Ref: https://stackoverflow.com/a/28466267/519360
    local _LONG_OPTARG=
    local _OPTARG=
    while getopts ':wvx:p:u:a:i:I:z2:-:' _OPTARG; do
        case "${_OPTARG}" in
            w | v)
                run_cmd 'true' ${_ARGS}
                break
                ;;
            x)
                halt 1 "${_OPTARG}"
                break
                ;;
            p)
                halt 2 "${_OPTARG}"
                break
                ;;
            u | a | i | I | z)
                halt 0 "${_OPTARG}"
                break
                ;;
            - )  _LONG_OPTARG="${OPTARG#*=}"
                case "${OPTARG}" in
                    help | version)
                        run_cmd 'true' ${_ARGS}
                        break
                        ;;
                    pid-file=?)
                        halt 1 "${OPTARG}"
                        break
                        ;;
                    port=?)
                        halt 2 "${OPTARG}"
                        break
                        ;;
                    user=? | listen-address=? | interface=? | except-interface=? | bind-interface | bind-dynamic | no-dhcp-interface=? | local-service)
                        halt 0 "${OPTARG}"
                        break
                        ;;
                esac ;;
            ?)
                # bypass the unknown flags to the software
                break
                ;;
        esac
    done
}

cmd_start() {
    local _ARGS="$@"

    # check the options
    check_opts ${_ARGS}

    # `-u root` is necessary to access NET_CAP_ADMIN in a container
    # even when the DHCP is disabled.
    # Ref: https://github.com/moby/moby/issues/4366
    # Additionally, pidfile is disabled for the read-only situation
    run_cmd 'true' \
            -p ${LISTEN_PORT} \
            -i eth0 \
            -2 eth0 -2 lo \
            --local-service \
            --pid-file -k \
            ${_ARGS}
}

cmd_sample() {
    cp -R "${SAMPLE_CFG_DIR}"/* "${CFG_DIR}"/
    info "A config sample is created.\nMake sure you have mounted a folder into ${CFG_DIR} before the container is created."
}

cmd_help() {
	cat <<- EOF
	Commands:
	    start (Default)     Start a server. You could load existing *.conf files by mounting them into the container at ${CFG_DIR}.
	    sample              Create a config sample in ${CFG_DIR}.
	    help                Show this help message.
	EOF
}

parse_cmd() {
    # debug
    #info "${FUNCNAME}" "$*"

    local _SUB_CMD="$1"
    local _ARGS="$@"

    case "${_SUB_CMD}" in
        start)
            shift
            _ARGS="$@"
            cmd_start ${_ARGS}
            ;;
        sample)
            cmd_sample
            ;;
        help)
            cmd_help
            ;;
        *)
            cmd_start ${_ARGS}
            ;;
    esac
}

main() {
    #info "${FUNCNAME}" "${ARGS}"

    # env check
    [[ -n "${LISTEN_PORT}" && -x "${ENTRYPOINT_CMD}" ]] \
        || fatal "This script is only compatible with the nutshells/dnsmasq-fast-lookup image,\ndo NOT try to run it out straight."

    parse_cmd ${ARGS}
}

trap cleanup EXIT
main
