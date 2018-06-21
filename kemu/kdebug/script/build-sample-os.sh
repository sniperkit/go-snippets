#!/bin/sh
if [ $# != 1 ]; then
    echo "Please give the build target!"
    exit 
fi

. $PROJECT/script/mkdev.sh
if [ $1 = "host1" ]; then
    rm -rf ../result/build/usr/initramfs_data.cpio.gz
    #cp -a ../samples/topology/start-host1 ../rootfs/etc/init.d/rcS
    cp -a ../../examples/rcS ../rootfs/etc/init.d/rcS
    cp -a ../../examples/busybox.inittab ../rootfs/etc/inittab
    cd .. &&  make kernel
elif [ $1 = "host2" ]; then
    rm -rf ../result/build/usr/initramfs_data.cpio.gz
    cp -a ../samples/topology/start-host2 ../rootfs/etc/init.d/rcS
    cd .. &&  make kernel
elif [ $1 = "router" ]; then
    rm -rf ../result/build/usr/initramfs_data.cpio.gz
    cp -a ../samples/topology/start-router ../rootfs/etc/init.d/rcS
    cd .. &&  make kernel
else
    echo "Unsupported target: $1 "
    exit -1
    # sudo qemu-system-x86_64 -m 2000 -net nic,model=e1000,vlan=2, -net tap,ifname=tap4,vlan=2 -kernel result/build/arch/x86/boot/bzImage -append vga=0x380 -vga vmware
fi

