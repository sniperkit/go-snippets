#!/bin/sh

if [ $# != 1 ]; then
    echo "Please give the build target!"
    exit 
fi

if [ "$ARCH"x = "arm"x ]; then
    sudo qemu-system-arm -M versatilepb -m 1024 -kernel result/build/arch/arm/boot/zImage  
elif [ "$ARCH"x = "x86_64"x ]; then
    if [ $1 = "host1" ]; then
        sudo  qemu-system-x86_64 -m 4096 -kernel result/build/arch/x86_64/boot/bzImage -append "console=ttyS0" -nographic
        #sudo qemu-system-x86_64 -m 2000 -net nic,model=e1000,vlan=1 --net tap,ifname=tap1,vlan=1,script=no,downscript=no -kernel result/build/arch/x86/boot/bzImage -append vga=0x380 -vga vmware &
    elif [ $1 = "host2" ]; then
        sudo qemu-system-x86_64 -m 2000 -net nic,model=e1000,vlan=2 -net tap,ifname=tap2,vlan=2,script=no,downscript=no -kernel result/build/arch/x86/boot/bzImage -append vga=0x380 -vga vmware &
    elif [ $1 = "router" ]; then
        sudo qemu-system-x86_64 -m 2000 -net nic,model=e1000,vlan=1 -net tap,ifname=tap3,vlan=1,script=no,downscript=no  -net nic,model=e1000,vlan=2 -net tap,ifname=tap4,vlan=2,script=no,downscript=no  -kernel result/build/arch/x86/boot/bzImage -append vga=0x380 -vga vmware &
    else
        echo "Unsupported target: $1 "
        exit -1
        # sudo qemu-system-x86_64 -m 2000 -net nic,model=e1000,vlan=2, -net tap,ifname=tap4,vlan=2 -kernel result/build/arch/x86/boot/bzImage -append vga=0x380 -vga vmware
    fi
else
    echo "Unsupported target platform"
    exit -1
fi
