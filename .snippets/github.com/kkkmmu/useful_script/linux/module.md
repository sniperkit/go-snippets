1. Makefile写法：
	obj-m += kprobe.o
	obj-m += jprobe.o
	obj-m += hw_breakpoint.o
2. 编译方法:
	 $(MAKE) -C $(KERNEL_DIR) M=`pwd` modules 
	 $(MAKE) -C $(KERNEL_DIR) M=`pwd` clean
3. 编译内核成功以后从http://fs.devloop.org.uk上下载rootfs. 然后按一下方式启动：
	sudo qemu-system-x86_64 -nographic -m 512 -kernel build/arch/x86_64/boot/bzImage -drive file=./CentOS6.x-AMD64-root_fs,if=ide -append "root=/dev/sda earlyprintk=serial,ttyS0,9600 console=ttyS0,9600n8"
