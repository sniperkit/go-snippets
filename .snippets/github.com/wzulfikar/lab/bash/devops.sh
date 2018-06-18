# service aliases
nginx_start="service nginx start"
nginx_stop="service nginx stop"

traccar_start="service traccar start"
traccar_stop="service traccar stop"

supervisor_start="service supervisor start"
supervisor_stop="service supervisor stop"

# logs aliases
# nginx_log="/var/log/nginx/error.log"

# log () {
#     run "cat $1"
# }

run () {
    LAST_ARG=${@: -1}
    if [ $LAST_ARG == "sudo" ]; then
        eval sudo $1
    else
        eval $1
    fi
}

start () {
    service=$1
    cmd=start
    str_exec=${service}_${cmd}
    echo "Running ${str_exec}.."
    run \$${service}_${cmd} ${@: -1}
}

stop () {
    service=$1
    cmd=stop
    str_exec=${service}_${cmd}
    echo "Stopping ${str_exec}.."
    run \$${service}_${cmd} ${@: -1}
}

restart () {
    stop $1
    start $1
}

# site util
site () {
    # paths
    SITE_PATH="~/desktop"
    SITE_TEMPLATE="~/desktop"

    # commands
    SITE_CREATE="touch ${SITE_PATH}/sites-available/$2"
    SITE_VIEW="cat ${SITE_PATH}/sites-enabled/$2"
    SITE_EDIT="nano ${SITE_PATH}/sites-available/$2"
    SITE_ENABLE="mv ${SITE_PATH}/sites-available/$2 ${SITE_PATH}/sites-enabled"
    SITE_UPDATE="ln -s ${SITE_PATH}/sites-available/$2 ${SITE_PATH}/sites-enabled"
    
    LAST_ARG=${@: -1}
    
    if [ $1 == "create" ]; then
        run ${SITE_CREATE} $LAST_ARG
    elif [ $1 == "enable" ]; then
        run ${SITE_ENABLE} $LAST_ARG
    elif [ $1 == "update" ]; then
        run ${SITE_UPDATE} $LAST_ARG
    else
        run ${1} $SITE_PATH $LAST_ARG
    fi
}
