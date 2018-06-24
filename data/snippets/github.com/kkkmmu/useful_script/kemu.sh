#/bin/sh

top=`pwd`
# Here we should use the source from network. .. Github
sudo dnf -y install vim make qemu bc cmake texinfo help2man patch libtool bison flex autoconf gperf
# sudo apt-get -y install vim make qemu bc

mkdir -p $top/tools
mkdir -p $top/tar

# For brctl and tunctl
sudo dnf -y install bridge-utils
sudo dnf -y install uml-utilities
sudo dnf -y install golang
sudo dnf -y install ssh
sudo dnf -y install git

#sudo apt-get -y install bridge-utils
#sudo apt-get -y install uml-utilities
#sudo apt-get -y install golang
#sudo apt-get -y install ssh
#sudo apt-get -y install git

cd $top/tools 
if [ ! -d crosstool-ng ]; then
    #
    # This is for build crosstool-ng
    # RUN apt-get -y install gperf bison flex texinfo help2man gawk automake
    # RUN git clone https://github.com/lamproae/lib-tool.git \
        # && cd lib-tool \
        # && ./bootstrap \
        # && ./configure \
        # && make \
        # && make install

    # This is for get croostool-ng
    # WORKDIR $HOME
    # git clone https://github.com/lamproae/crosstool-ng.git \
    git clone https://github.com/crosstool-ng/crosstool-ng.git \
        && cd $top/tools/crosstool-ng \
        && ./bootstrap \
        && ./configure \
        && make \
        && sudo make install

    # Build cross-compile toolchains
    ct-ng x86_64-unknown-linux-gnu \
        && ct-ng build
fi

cd $top/ 
if [ ! -d kdebug ]; then
    git clone https://github.com/lamproae/kdebug.git
fi

cd $top/ 
if [ ! -d kdebug/tools/toolchain/x86_64 ]; then
    mkdir -p $top/kdebug/tools/toolchain/x86_64
    cp -a $HOME/x-tools/x86_64-unknown-linux-gnu $top/kdebug/tools/toolchain/x86_64
fi

# Get kernel source code
cd $top/tar 
if [ ! -d linux-4.10.12 ]; then
    curl -o linux-4.10.12.tar.xz https://www.kernel.org/pub/linux/kernel/v4.x/linux-4.10.12.tar.xz \
        && xz -d linux-4.10.12.tar.xz \
        && tar -xvf linux-4.10.12.tar 
    cd ..
    mkdir -p $top/kdebug/kernel
    cp -a tar/linux-4.10.12 $top/kdebug/kernel/
fi

# Get busybox source code
cd $top/tar 
if [ ! -d busybox-1.24.2 ]; then
    curl -o busybox-1.24.2.tar.bz2 https://www.busybox.net/downloads/busybox-1.24.2.tar.bz2 \
        && tar -jxvf busybox-1.24.2.tar.bz2 
    cd ..
    mkdir -p $top/kdebug/apps/busybox
    cp -a tar/busybox-1.24.2 $top/kdebug/apps/busybox/
fi

# Get coreutils source code
cd $top/tar 
if [ ! -d coreutils-8.23 ]; then
    curl -o coreutils-8.23.tar.xz http://ftp.gnu.org/gnu/coreutils/coreutils-8.23.tar.xz\
        && xz -d coreutils-8.23.tar.xz \
        && tar -xvf coreutils-8.23.tar 
    cd ..
    mkdir -p $top/kdebug/apps/coreutils
    cp -a tar/coreutils-8.23 $top/kdebug/apps/coreutils/
fi

# Get inetutils source code
cd $top/tar 
if [ ! -d inetutils-1.9.4 ]; then
    curl -o inetutils-1.9.4.tar.xz http://ftp.gnu.org/gnu/inetutils/inetutils-1.9.4.tar.xz \
        && xz -d inetutils-1.9.4.tar.xz \
        && tar -xvf inetutils-1.9.4.tar 
    cd ..
    mkdir -p $top/kdebug/apps/inetutils
    cp -a tar/inetutils-1.9.4 $top/kdebug/apps/inetutils/
fi

cd $top/tar 
if [ ! -d bash-4.3 ]; then
    curl -o bash-4.3.tar.gz http://ftp.gnu.org/gnu/bash/bash-4.3.tar.gz \
        && tar -zxvf bash-4.3.tar.gz 
    cd ..
    mkdir -p $top/kdebug/apps/bash
    cp -a tar/bash-4.3 $top/kdebug/apps/bash/
fi

cd $top/tar 
if [ ! -d ncurses-5.4 ]; then
    curl -o ncurses-5.4.tar.gz http://ftp.gnu.org/gnu/ncurses/ncurses-5.4.tar.gz \
        && tar -zxvf ncurses-5.4.tar.gz 
    cd ..
    mkdir -p $top/kdebug/apps/ncurses
    cp -a tar/ncurses-5.4 $top/kdebug/apps/ncurses/
fi

cd $top/tar 
if [ ! -d ncurses-5.4 ]; then
    curl -o ncurses-5.4.tar.gz http://ftp.gnu.org/gnu/ncurses/ncurses-5.4.tar.gz \
        && tar -zxvf ncurses-5.4.tar.gz 
    cd ..
    mkdir -p $top/kdebug/apps/ncurses
    cp -a tar/ncurses-5.4 $top/kdebug/apps/ncurses/
fi

cd $top/tar 
if [ ! -d iptables-1.4.20 ]; then
    curl -o iptables-1.4.20.tar.bz2 http://www.netfilter.org/projects/iptables/files/iptables-1.4.20.tar.bz2 \
        && tar -jxvf iptables-1.4.20.tar.bz2 
    cd ..
    mkdir -p $top/kdebug/apps/iptables
    cp -a tar/iptables-1.4.20 $top/kdebug/apps/iptables
fi

cd $top/tar 
if [ ! -d procps-3.2.8 ]; then
    git clone https://github.com/lamproae/procps-3.2.8.git
    cd ..
    mkdir -p $top/kdebug/apps/procps
    cp -a tar/procps-3.2.8 $top/kdebug/apps/procps/
fi

cd $top/tar 
if [ ! -d iproute2 ]; then
    git clone https://github.com/lamproae/iproute2.git
    cd ..
    cp -a tar/iproute2 $top/kdebug/apps/iproute2
fi

cd $top/tar 
if [ ! -d tcpdump-4.7.4 ]; then 
    curl -o tcpdump-4.7.4.tar.gz http://www.tcpdump.org/release/tcpdump-4.7.4.tar.gz \
        && tar -zxvf tcpdump-4.7.4.tar.gz 
    curl -o libpcap-1.6.2.tar.gz http://www.tcpdump.org/release/libpcap-1.6.2.tar.gz \
        && tar -zxvf libpcap-1.6.2.tar.gz
    cd ..
    mkdir -p $top/kdebug/apps/tcpdump
    cp -a tar/tcpdump-4.7.4 $top/kdebug/apps/tcpdump/
    cp -a tar/libpcap-1.6.2 $top/kdebug/apps/tcpdump/
fi

cd $top/tar 
if [ ! -d strace-4.8 ]; then
    curl -o strace-4.8.tar.xz http://pkgs.fedoraproject.org/repo/pkgs/strace/strace-4.8.tar.xz/c575ef43829586801f514fd91bfe7575/strace-4.8.tar.xz \
        && xz -d strace-4.8.tar.xz && tar -xvf strace-4.8.tar 
    cd ..
    mkdir -p $top/kdebug/apps/strace
    cp -a tar/strace-4.8 $top/kdebug/apps/strace/
fi

cd $top/tar 
if [ ! -d quagga ]; then 
    git clone https://github.com/lamproae/quagga.git
    cd ..
    mkdir -p $top/kdebug/apps/quagga
    cp -a tar/quagga $top/kdebug/apps/quagga
fi

cd $top/tar 
if [ ! -d gawk-4.1.0 ]; then
    curl -o gawk-4.1.0.tar.xz http://ftp.gnu.org/gnu/gawk/gawk-4.1.0.tar.xz \
        && xz -d gawk-4.1.0.tar.xz && tar -xvf gawk-4.1.0.tar
    cd ..
    mkdir -p $top/kdebug/apps/gawk
    cp -a tar/gawk-4.1.0 $top/kdebug/apps/gawk/
fi

cd $top
# RUN apt-get update
# RUN apt-get -y install vim make
