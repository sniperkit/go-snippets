#!/bin/sh

if [ $# != 1 ];then
    echo "Please give the operation!"
    exit -1 
fi

if [ $1 = "set" ];then

    sudo ifconfig eth0 down

    # Add bridge 0
    sudo brctl addbr br0
    sudo brctl addif br0 eth0
    sudo brctl stp br0 off
    sudo brctl setfd br0 1
    sudo brctl sethello br0 1

    sudo ifconfig br0 0.0.0.0 promisc up
    sudo ifconfig eth0 0.0.0.0 promisc up

    sudo ifconfig br0 10.71.1.193 netmask 255.255.255.0
    sudo route add -net 0.0.0.0 netmask 0.0.0.0 gw 10.71.1.254

    sudo brctl show br0
    sudo brctl showstp br0

    sudo tunctl -t tap1 -u root
    sudo tunctl -t tap2 -u root
    sudo tunctl -t tap3 -u root
    sudo tunctl -t tap4 -u root
    sudo tunctl -t tap5 -u root
    sudo brctl addif br0 tap1
    sudo brctl addif br0 tap2
    sudo brctl addif br0 tap3
    sudo brctl addif br0 tap4
    sudo brctl addif br0 tap5
    sudo ifconfig tap1 0.0.0.0 promisc up
    sudo ifconfig tap2 0.0.0.0 promisc up
    sudo ifconfig tap3 0.0.0.0 promisc up
    sudo ifconfig tap4 0.0.0.0 promisc up
    sudo ifconfig tap5 0.0.0.0 promisc up
    sudo brctl showstp br0
elif [ $1 = "reset" ];then
    sudo ifconfig eth0 down
    sudo brctl delif br0 eth0
    sudo ifconfig br0 down
    sudo ifconfig tap1 down
    sudo ifconfig tap2 down
    sudo ifconfig tap3 down
    sudo ifconfig tap4 down
    sudo ifconfig tap5 down
    sudo brctl delbr br0
    sleep 2
    sudo ifconfig eth0 up
fi
