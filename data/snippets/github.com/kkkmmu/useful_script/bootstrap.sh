#!/bin/bash

top=`pwd`

# shell script dictionary
declare -A Download
declare -A Decompress
declare -A Compress

Compress=(["tgz"]="tar.gz" ["txz"]="tar.xz" ["zip"]="gzip" ["tar"]="tar" ["zip2"]="tar.bz2")
Download=(["git"]="git clone " ["wget"]="wget ")
Decompress=(["tgz"]="tar -zxvf " ["txz"]="xz -d " ["zip"]="unzip " ["tar"]="tar -xvf " ["zip2"]="tar -jxvf ")

function process() 
{
    i=0
    while [ $i -lt 1000 ]
    do
        for j in '-' '\' '|' '/'
        do
            printf "                %s\r" $j
            sleep 0.5
            ((i++))
        done
    done
}

function execute()
{
    #$* 
    $* &>> $top/kemu.log
    if [ $? != 0 ]; then
        echo "Command \"$*\" excute failed!";
        exit 1;
    fi
}

function download_package()
{
    cd $top/tar
    echo " Downloading Package $pname !"
    execute ${Download[$1]} $2
    echo " Package $pname Downloading finished.!"
    cd $top
}

function decompress_package()
{
    cd $top/tar
    echo " Decompressing Package $pname !"
    execute ${Decompress[$2]} $1.${Compress[$2]}
    if [ $2 == "txz" ]; then
        execute ${Decompress["tar"]} $1.${Compress["tar"]};
    fi
    echo " Package $pname Decompressing finished!"
    cd $top
}

function locallize_package()
{
    cd $top/tar
    echo " Locallizing Package $pname !"
    #if [ $1 != $2 ];then
        execute mv $1 $top/kdebug/apps/$2
    #fi
    echo " Package $pname Locallizing finished!"
    cd $top
}

function _prepare_package() 
{
    local pname=$1;
    local pdown=$2;
    local pcmpz=$3;
    local lname=$4;
    local plink=$5;
    download_package $pdown $plink
    if [ $pdown != "git" ]; then
        decompress_package $pname $pcmpz;
    fi
    locallize_package $pname $lname;
    #TODO: If failed need to clean
    if [ $? == 0 ]; then
        echo 1 >> .$lname.p
    fi 
}

function prepare_package()
{
    cat .package | while read p; do
    {
        pname=`echo $p | awk '{split($0, a," ");print a[1]}'`;
        pvers=`echo $p | awk '{split($0, a," ");print a[2]}'`;
        pdown=`echo $p | awk '{split($0, a," ");print a[3]}'`;
        pcmpz=`echo $p | awk '{split($0, a," ");print a[4]}'`;
        lname=`echo $p | awk '{split($0, a," ");print a[5]}'`;
        plink=`echo $p | awk '{split($0, a," ");print a[6]}'`;

        if [ $pvers != "-" ]; then
            pname=$pname-$pvers
        fi

        if [ -f .$lname.p ]; then
            continue
        else
            clean $lname
        fi

        echo "==================================================="
        echo "    Processing Package $pname        "
        _prepare_package $pname $pdown $pcmpz $lname $plink
        echo "    Package $pname process finished   "
        echo "==================================================="
    };
done
}

function clean() 
{
    cat .package | while read p; do
    {
        pname=`echo $p | awk '{split($0, a," ");print a[1]}'`;
        pvers=`echo $p | awk '{split($0, a," ");print a[2]}'`;
        pdown=`echo $p | awk '{split($0, a," ");print a[3]}'`;
        pcmpz=`echo $p | awk '{split($0, a," ");print a[4]}'`;
        lname=`echo $p | awk '{split($0, a," ");print a[5]}'`;
        plink=`echo $p | awk '{split($0, a," ");print a[6]}'`;

        if [ $# == 1 ]; then
            if [ $lname == $1 ]; then
                exist=`ls -al | grep $pname`
                if [ '$exist'x != ''x ]; then
                    rm -rf $lname*;
                fi

                if [ -f .$lname.p ]; then
                    rm -rf .$lname.p;
                fi
            else
                continue
            fi
        else
            exist=`ls -al | grep $pname`
            if [ '$exist'x != ''x ]; then
                rm -rf $lname*;
            fi

            if [ -f .$lname.p ]; then
                rm -rf .$lname.p;
            fi

            if [ -f .env.p ]; then
                rm -rf .env.p
            fi

            if [ -f .crosstool-ng.p ]; then
                rm -rf .crosstool-ng.p
            fi
        fi
    };
done
}


function prepare_env() 
{
    sudo dnf -y install vim make qemu bc cmake texinfo help2man patch libtool bison flex autoconf gperf tunctl
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
    mkdir -p $top/tar
    mkdir -p $top/tools
    mkdir -p $top/kdebug/tools/toolchain
    mkdir -p $top/kdebug/kernel

    echo 1 >> .env.p
}

function prepare_toolchain() 
{
    if [ ! -d crosstool-ng ]; then
        git clone https://github.com/crosstool-ng/crosstool-ng.git 
    fi

    cd $top/tools/crosstool-ng \
        && ./bootstrap \
        && ./configure \
        && make \
        && sudo make install

    # Build cross-compile toolchains
    ct-ng x86_64-unknown-linux-gnu \
        && ct-ng CT_PREFIX=$top build

    sudo cp -a $top/x86_64-unknown-linux-gnu $top/kdebug/tools/toolchain/
    echo 1 >> .crosstool-ng.p
}


function prepare_misc()
{
    if [ -d $top/kdebug/apps/libpcap ]; then
        mv $top/kdebug/apps/libpcap $top/kdebug/apps/tcpdump/libpcap
    fi

    if [ -d $top/kdebug/apps/linux ]; then
        mv $top/kdebug/apps/linux $top/kdebug/kernel/
    fi

    export PROJECT=$top/kdebug/
    echo 1 >> .misc.p
}

function prepare()
{
    if [ ! -f .env.p ]; then
        prepare_env
    fi

    if [ ! -f .crosstool-ng.p ]; then
        prepare_toolchain
    fi

    prepare_package;

    if [ ! -f .misc.p ]; then
        prepare_misc
    fi
}
prepare
