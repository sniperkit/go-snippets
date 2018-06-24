#!/bin/sh

devdir=$PROJECT/rootfs/dev

function mkdev() {
    echo "Creating Initial device node"
    if [ ! -f $devdir/console ]; then
        if [ ! -d $devdir ]; then
            mkdir -p $devdir
        fi
        sudo mknod $devdir/console c 5 1
        sudo mknod $devdir/full c 1 7
        sudo mknod $devdir/kmem c 1 2
        sudo mknod $devdir/mem c 1 1
        sudo mknod $devdir/null c 1 3
        sudo mknod $devdir/port c 1 4
        sudo mknod $devdir/random c 1 8
        sudo mknod $devdir/urandom c 1 9
        sudo mknod $devdir/zero c 1 5
    fi
}

mkdev
