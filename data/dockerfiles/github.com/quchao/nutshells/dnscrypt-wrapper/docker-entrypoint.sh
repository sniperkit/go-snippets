#!/bin/sh
set -eo pipefail

# constants
readonly ARGS="$@"
readonly ENTRYPOINT_CMD='/usr/local/sbin/dnscrypt-wrapper'
readonly KEYS_DIR="${KEYS_DIR:-/usr/local/etc/dnscrypt-wrapper}"
readonly CRYPT_KEYS_DIR="${CRYPT_KEYS_DIR:-${KEYS_DIR}/crypt}"
readonly PUB_KEY_FILENAME="public.key"
readonly PRI_KEY_FILENAME="secret.key"
readonly PUB_KEY="${KEYS_DIR}/${PUB_KEY_FILENAME}"
readonly PRI_KEY="${KEYS_DIR}/${PRI_KEY_FILENAME}"

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
            _ERR_MSG="${_ERR_MSG}Don't try to fork it into background in a container.\n"
            ;;
        2)
            # built in
            _ERR_MSG="${_ERR_MSG}Use the '-e', '--env' or '--env-file' options to
            override its value,\n"
            _ERR_MSG="${_ERR_MSG}referer to README for the relevant environment variables.\n"
            ;;
    esac

    _ERR_MSG="${_ERR_MSG}Just simply remove it and "

    case "${_ERR_NUM}" in
        3)
            # run sub-cmd
            local _SUB_CMD="${3:-start}"
            _ERR_MSG="${_ERR_MSG}run the '${_SUB_CMD}' command instead."
            ;;
        *)
            # default
            _ERR_MSG="${_ERR_MSG}try again."
            ;;
    esac

    fatal "${_ERR_MSG}"
}

run_cmd() {
    local _EXEC="$1"
    shift
    local _ARGS="$@"

    if [[ "${_EXEC}" == 'true' ]]; then
        exec "${ENTRYPOINT_CMD}" \
            -u "${RUN_AS_USER}" \
            ${_ARGS}
    else
        "${ENTRYPOINT_CMD}" \
            -u "${RUN_AS_USER}" \
            ${_ARGS}
    fi
}

check_opts() {
    local _ARGS="$@"

    # Ref: https://stackoverflow.com/a/28466267/519360
    local _LONG_OPTARG=
    local _OPTARG=
    while getopts ':hvdp:u:a:r:-:o:x' _OPTARG; do
        case "${_OPTARG}" in
            h | v)
                run_cmd 'true' ${_ARGS}
                break
                ;;
            d | p)
                halt 1 "${_OPTARG}"
                break
                ;;
            u | a | r)
                halt 2 "${_OPTARG}"
                break
                ;;
            - )  _LONG_OPTARG="${OPTARG#*=}"
                case "${OPTARG}" in
                    help | version)
                        run_cmd 'true' ${_ARGS}
                        break
                        ;;
                    daemonize | pidfile=?)
                        halt 1 "${OPTARG}"
                        break
                        ;;
                    user=? | listen-address=? | resolver-address=? | provider-name=?)
                        halt 2 "${OPTARG}"
                        break
                        ;;
                    gen-provider-keypair)
                        halt 3 "${OPTARG}" 'init'
                        break
                        ;;
                    show-provider-publickey)
                        halt 3 "${OPTARG}" 'pubkey'
                        break
                        ;;
                    show-provider-publickey-dns-records)
                        halt 3 "${OPTARG}" 'dns'
                        break
                        ;;
                    outgoing-address=? | xchacha20 | provider-cert-file=? | crypt-secretkey-file=? | gen-cert-file | gen-crypt-keypair | provider-publickey-file=? | provider-secretkey-file=? | cert-file-expire-days=?)
                        halt 0 "${OPTARG}"
                        break
                        ;;
                esac ;;
            o | x)
                halt 0 "${_OPTARG}"
                break
                ;;
            ?)
                # bypass the unknown flags to the software
                break
                ;;
        esac
    done
}

is_initialized() {
    if [[ ! -f "${PUB_KEY}" || ! -f "${PRI_KEY}" ]]; then
        warning 'The provider key pair does NOT exist.'
        return 1
    else
        return 0
    fi
}

ensure_initialized() {
    is_initialized \
        || fatal "$(cat <<- EOF
		Not initialized yet.
		Run the 'init' command to generate a new provider key pair,
		or use existing ones by mounting them into the container at ${KEYS_DIR}.
		Referer to README for more details please.
	EOF
	)"
}

need_rotation() {
    local _LIFESPAN="$((CRYPT_KEYS_LIFESPAN * 1440 * 7 / 10))"
    if [[ \
        "$(/usr/bin/find "${CRYPT_KEYS_DIR}" -maxdepth 1 -type f \
            -name '*.key' -mmin -"${_LIFESPAN}" \
            -print | wc -l | sed 's|[^0-9]||g')" -eq 0 \
         ]]; then
        warning "Rotation is needed for no crypt key is valid until ${_LIFESPAN} minutes later."
        return 0
    else
        return 1
    fi
}

prune_keys() {
    # prune the expired keys/certs
    /usr/bin/find "${CRYPT_KEYS_DIR}" -maxdepth 1 -type f \
        -mmin +"$((CRYPT_KEYS_LIFESPAN * 1440))" \
        -exec rm -f {} \;
}

rotate_keys() {
    local _TS
    _TS="$(date '+%s')"
    local _CRYPT_KEY="${CRYPT_KEYS_DIR}/${_TS}.key"
    local _CRYPT_XSALSA20_CERT="${CRYPT_KEYS_DIR}/${_TS}-xsalsa20.cert"
    local _CRYPT_XCHACHA20_CERT="${CRYPT_KEYS_DIR}/${_TS}-xchacha20.cert"

    # generate a key pair for encryption, and sign 2 certs
    run_cmd 'false' --gen-crypt-keypair \
            --crypt-secretkey-file="${_CRYPT_KEY}" \
        && run_cmd 'false' --gen-cert-file \
            --provider-publickey-file="${PUB_KEY}" \
            --provider-secretkey-file="${PRI_KEY}" \
            --crypt-secretkey-file="${_CRYPT_KEY}" \
            --provider-cert-file="${_CRYPT_XSALSA20_CERT}" \
            --cert-file-expire-days="${CRYPT_KEYS_LIFESPAN}" \
        && run_cmd 'false' --gen-cert-file \
            --provider-publickey-file="${PUB_KEY}" \
            --provider-secretkey-file="${PRI_KEY}" \
            --crypt-secretkey-file="${_CRYPT_KEY}" \
            --provider-cert-file="${_CRYPT_XCHACHA20_CERT}" \
            --cert-file-expire-days="${CRYPT_KEYS_LIFESPAN}" \
            --xchacha20

    # set permissions
    chmod 644 "${_CRYPT_XSALSA20_CERT}" "${_CRYPT_XCHACHA20_CERT}"
    chmod 640 "${_CRYPT_KEY}"
    chgrp "${RUN_AS_USER}" "${_CRYPT_KEY}" "${_CRYPT_XSALSA20_CERT}" "${_CRYPT_XCHACHA20_CERT}"

    info 'Crypt key has been rotated.'
}


list_crypt_key_pair() {
    # list 2 most recent modified crypt keys/certs
    local _EXT="${1:-key}"
    local _NUM="${2:-2}"

    local _FILE=
    local _FILE_LIST=
    for _FILE in $(/usr/bin/find "${CRYPT_KEYS_DIR}" -maxdepth 1 -type f -name "*.${_EXT}" -print | sort -nr | head -"${_NUM}"); do
        _FILE_LIST="${_FILE_LIST}${_FILE},"
    done

    echo "${_FILE_LIST}"
}

cmd_start() {
    local _ARGS="$@"

    # check the options
    check_opts ${_ARGS}

    # ensure it's initialized
    ensure_initialized

    # prepare for crypt key & certs
    mkdir -p "${CRYPT_KEYS_DIR}"

    # prune the oldies
    prune_keys

    # check if the crypt key need a rotation
    need_rotation \
        && rotate_keys

    # get ip of the linked container
    if echo "${RESOLVER_IP}" | grep -Eqv '^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$'; then
        local _RESOLVER_IP
        set +e
        _RESOLVER_IP="$(getent hosts "${RESOLVER_IP}" | awk '{ print $1 }')"
        set -e
        if [[ -n "${_RESOLVER_IP}" ]]; then
            info "'${RESOLVER_IP}' is resolved to ${_RESOLVER_IP}."
            RESOLVER_IP="${_RESOLVER_IP}"
        fi
    fi

    run_cmd 'true' \
            -a "0.0.0.0:${LISTEN_PORT}" \
            -r "${RESOLVER_IP}:${RESOLVER_PORT}" \
            --provider-name="2.dnscrypt-cert.${PROVIDER_BASENAME}" \
            --provider-cert-file="$(list_crypt_key_pair 'cert' '4')" \
            --crypt-secretkey-file="$(list_crypt_key_pair)" \
            ${_ARGS}
}

cmd_init() {
    local _ARGS="$@"

    # initialized already
    if is_initialized ; then
        info 'The container has been initialized, starting now...'
        cmd_start ${_ARGS}
        exit $?
    fi

    # generate provider key pair
    run_cmd 'false' --gen-provider-keypair \
         && mv "${PUB_KEY_FILENAME}" "${PRI_KEY_FILENAME}" "${KEYS_DIR}"
    #chmod 644 "${PUB_KEY}"
    #chmod 640 "${PRI_KEY}"

    # recheck
    if is_initialized ; then
        cmd_start ${_ARGS}
    else
        fatal 'Failed to initialize the container, file an
        issue on GitHub please.'
    fi
}

cmd_pubkey() {
    # ensure it's initialized
    ensure_initialized

    run_cmd 'true' \
        --show-provider-publickey \
        --provider-publickey-file="${PUB_KEY}"
}

cmd_dns() {
    # ensure it's initialized
    ensure_initialized

    # ensure the crypt key correctly generated
    [[ -e "${CRYPT_KEYS_DIR}" ]] \
        || fatal "Crypt key not found, run the 'start' command to
        generate one please."

    run_cmd 'true' \
        --show-provider-publickey-dns-records \
        --provider-cert-file="$(list_crypt_key_pair 'cert' '4')"
}

cmd_help() {
	cat <<- "EOF"
	Commands:
	    init                Perform an initialization, which technically generates a new provider key pair before starting the server.
	    start (Default)     Start a server. It will generate the crypt key & certs, and rotate them every time it starts if necessary.
	    pubkey              Show public key's fingerprint.
	    dns                 Show DNS record.
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
        init)
            shift
            _ARGS="$@"
            cmd_init ${_ARGS}
            ;;
        pubkey)
            cmd_pubkey
            ;;
        dns)
            cmd_dns
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
    [[ -n "${LISTEN_PORT}" && -n "${RUN_AS_USER}" && -x "${ENTRYPOINT_CMD}" ]] \
        || fatal 'This script is only compatible with the nutshells/dnscrypt-wrapper image,\ndo NOT try to run it out straight.'

    parse_cmd ${ARGS}
}

trap cleanup EXIT
main
