#!/bin/sh -e

export OPENLDAP_ETC_DIR="/etc/openldap"
export OPENLDAP_RUN_DIR="/var/run/openldap"
export OPENLDAP_RUN_ARGSFILE="${OPENLDAP_RUN_DIR}/slapd.args"
export OPENLDAP_RUN_PIDFILE="${OPENLDAP_RUN_DIR}/slapd.pid"
export OPENLDAP_MODULES_DIR="/usr/lib/openldap"
export OPENLDAP_CONFIG_DIR="${OPENLDAP_ETC_DIR}/slapd.d"
export OPENLDAP_BACKEND_DIR="/var/lib/openldap"
export OPENLDAP_BACKEND_DATABASE="hdb"
export OPENLDAP_BACKEND_OBJECTCLASS="olcHdbConfig"
export OPENLDAP_ULIMIT="2048"

export LDAP_DOMAIN="${LDAP_DOMAIN:-example.com}"
export LDAP_DOMAIN_RDC="$(echo ${LDAP_DOMAIN} | sed 's/^\.//; s/\..*$//')"
# export LDAP_DOMAIN_OBJECTCLASS="organization
# objectClass: dcObject
# "
export LDAP_DOMAIN_OBJECTCLASS="${LDAP_DOMAIN_OBJECTCLASS:-domain}"
export LDAP_ORGANIZATION="${LDAP_ORGANIZATION:-${LDAP_DOMAIN}}"
export LDAP_SUFFIX="${LDAP_SUFFIX:-$(echo dc=$(echo ${LDAP_DOMAIN} | sed 's/^\.//; s/\./,dc=/g'))}"
export LDAP_PASSWORD="${LDAP_PASSWORD:-lderp!${LDAP_DOMAIN_RDC}}"
export LDAP_PASSWORD_ENCRYPTED="$(slappasswd -u -h '{SSHA}' -s ${LDAP_PASSWORD})"

ulimit -n ${OPENLDAP_ULIMIT}

if [[ ! -d ${OPENLDAP_CONFIG_DIR}/cn=config ]]; then
    mkdir -p ${OPENLDAP_CONFIG_DIR}

    if [[ ! -s ${OPENLDAP_ETC_DIR}/slapd-config.ldif ]]; then
        cat /srv/openldap/slapd-config.ldif.template | envsubst > ${OPENLDAP_ETC_DIR}/slapd-config.ldif
    fi

    slapadd -n0 -F ${OPENLDAP_CONFIG_DIR} -l ${OPENLDAP_ETC_DIR}/slapd-config.ldif > ${OPENLDAP_ETC_DIR}/slapd-config.ldif.log

    if [[ ! -s ${OPENLDAP_ETC_DIR}/ldap.conf ]]; then
        cat /srv/openldap/ldap.conf.template | envsubst > ${OPENLDAP_ETC_DIR}/ldap.conf
    fi

    mkdir -p ${OPENLDAP_BACKEND_DIR}/run
    chown -R ldap:ldap ${OPENLDAP_BACKEND_DIR}
    chown -R ldap:ldap ${OPENLDAP_CONFIG_DIR} ${OPENLDAP_BACKEND_DIR}

    if [[ -d /srv/openldap.d ]]; then
        if [[ ! -s /srv/openldap.d/000-domain.ldif ]]; then
            cat /srv/openldap/domain.ldif.template | envsubst > /srv/openldap.d/000-domain.ldif
        fi

        slapd_exe=$(which slapd)
        echo >&2 "$0 ($slapd_exe): starting initdb daemon"
        slapd -u ldap -g ldap -h ldapi:///

        for f in $(find /srv/openldap.d -type f | sort); do
            case "$f" in
                *.sh)   echo "$0: sourcing $f"; . "$f" ;;
                *.ldif) echo "$0: applying $f"; ldapadd -Y EXTERNAL -f "$f" 2>&1;;
                *)      echo "$0: ignoring $f" ;;
            esac
        done

        if [[ ! -s ${OPENLDAP_RUN_PIDFILE} ]]; then
            echo >&2 "$0 ($slapd_exe): ${OPENLDAP_RUN_PIDFILE} is missing, did the daemon start?"
            exit 1
        else
            slapd_pid=$(cat ${OPENLDAP_RUN_PIDFILE})
            echo >&2 "$0 ($slapd_exe): sending SIGINT to initdb daemon with pid=$slapd_pid"
            kill -s INT "$slapd_pid" || true
            while : ; do
                [[ ! -f ${OPENLDAP_RUN_PIDFILE} ]] && break
                sleep 1
                echo >&2 "$0 ($slapd_exe): initdb daemon is still up, sleeping ..."
            done
            echo >&2 "$0 ($slapd_exe): initdb daemon stopped"
        fi
    fi
fi

exec "$@"
