#!/bin/sh
set -eo pipefail

# constants
readonly ARGS="$@"

# global
SS_MODE='server'

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
        3)
            # unsupported
            _ERR_MSG="${_ERR_MSG}The manager/tunnel/redir modes haven't been supported yet by this image.\n"
            ;;
    esac

    _ERR_MSG="${_ERR_MSG}Just simply remove it and try again."

    fatal "${_ERR_MSG}"
}

run_cmd() {
    local _EXEC="$1"
    shift
    local _ARGS="$@"
    local _CMD="/usr/local/bin/ss-${SS_MODE}"

    if [[ "${_EXEC}" == 'true' ]]; then
        exec su-exec \
            "${RUN_AS_USER}" \
             "${_CMD}" \
            ${_ARGS}
    else
        "${_CMD}" \
            -a "${RUN_AS_USER}" \
            ${_ARGS}
    fi
}

check_opts() {
    local _ARGS="$@"

    # Ref: https://stackoverflow.com/a/28466267/519360
    local _LONG_OPTARG=
    local _OPTARG=
    while getopts ':hf:s:p:l:k:m:uUL:a:i:b:-:' _OPTARG; do
        case "${_OPTARG}" in
            h)
                run_cmd 'true' ${_ARGS}
                break
                ;;
            f)
                halt 1 "${_OPTARG}"
                break
                ;;
            s | p | l | k | m | u | U)
                halt 2 "${_OPTARG}"
                break
                ;;
            L)
                halt 3 "${_OPTARG}"
                break
                ;;
            a | i | b)
                halt 0 "${_OPTARG}"
                break
                ;;
            - )  _LONG_OPTARG="${OPTARG#*=}"
                case "${OPTARG}" in
                    help)
                        run_cmd 'true' ${_ARGS}
                        break
                        ;;
                    pid-file=?)
                        halt 1 "${OPTARG}"
                        break
                        ;;
                    server_host=? | server_port=? | local_port=? | password=? | key=? | encrypt_method=? | plugin=? | plugin-opts=?)
                        halt 2 "${OPTARG}"
                        break
                        ;;
                    manager-address=? | executable=?)
                        halt 3 "${OPTARG}"
                        break
                        ;;
                    user=? | interface=? | local_address=?)
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
    local _OBFS_OPTS=''

    # check the options
    check_opts ${_ARGS}

    # running mode
    if [[ "${SS_MODE}" == 'server' ]]; then
        _ARGS="${_ARGS} -s 0.0.0.0 -p ${LISTEN_PORT}"
    else
        [[ -n "${SERVER_ADDRESS}" || -n "${SERVER_PORT}" ]] \
            || fatal "Set 'SERVER_ADDRESS' and 'SERVER_PORT' both please."
        _ARGS="${_ARGS} -s ${SERVER_ADDRESS} -p ${SERVER_PORT} -b 0.0.0.0 -l ${LISTEN_PORT}"
    fi

    # password
    if [[ -z "${PASSWORD}" && -z "${KEY_IN_BASE64}" ]]; then
        fatal "Set either 'PASSWORD' or 'KEY_IN_BASE64' please."
    elif [[ -n "${PASSWORD}" ]]; then
        _ARGS="${_ARGS} -k ${PASSWORD}"
    else
        _ARGS="${_ARGS} --key ${KEY_IN_BASE64}"
    fi

    # other configs
    [[ "${REUSE_PORT}" == 'true' ]] \
        && _ARGS="${_ARGS} --reuse-port"
    [[ "${TCP_FAST_OPEN}" == 'true' ]] \
        && _ARGS="${_ARGS} --fast-open" \
        && _OBFS_OPTS=';fast-open'
    if [[ "${UDP_RELAY}" == 'true' ]]; then
        if [[ "${TCP_RELAY}" != 'true' ]]; then
            [[ "${SS_MODE}" != 'server' ]] \
                || fatal "'TCP_RELAY' must be enabled in local/client mode."
            _ARGS="${_ARGS} -U"
        else
            _ARGS="${_ARGS} -u"
        fi
    fi

    # obfs plugin
    if [[ "${WITH_OBFS}" == 'true' ]]; then
        [[ "${OBFS_PLUGIN}" != 'http' && "${OBFS_PLUGIN}" != 'tls' ]] \
            && fatal "'OBFS_PLUGIN' accepts either 'http' or 'tls' for now."
        if [[ "${SS_MODE}" == 'server' ]]; then
            _ARGS="${_ARGS} --plugin obfs-${SS_MODE} --plugin-opts obfs=${OBFS_PLUGIN}${_OBFS_OPTS}"
        else
            _ARGS="${_ARGS} --plugin obfs-${SS_MODE} --plugin-opts obfs=${OBFS_PLUGIN};obfs-host=${OBFS_HOST}${_OBFS_OPTS}"
        fi
    fi

    run_cmd 'true' \
            -m "${ENCRYPT_METHOD}" \
            ${_ARGS}
}

cmd_help() {
	cat <<- EOF
	Commands:
	    server (Default)    Start a server.
	    client              Start a client.
	    help                Show this help message.
	EOF
}

parse_cmd() {
    # debug
    #info "${FUNCNAME}" "$*"

    local _SUB_CMD="$1"
    local _ARGS="$@"

    case "${_SUB_CMD}" in
        server | start)
            shift
            _ARGS="$@"
            cmd_start ${_ARGS}
            ;;
        client | local)
            shift
            _ARGS="$@"
            SS_MODE='local'
            cmd_start ${_ARGS}
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
    [[ -n "${LISTEN_PORT}" && -n "${RUN_AS_USER}" ]] \
        || fatal 'This script is only compatible with the nutshells/ss-obfs image,\ndo NOT try to run it out straight.'

    parse_cmd ${ARGS}
}

trap cleanup EXIT
main
