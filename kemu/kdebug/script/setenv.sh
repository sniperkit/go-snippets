#!/bin/sh
function mithrandir()
{
    local target

    local uname=$(uname -msr)
    echo
    echo "Your are buiding on: $uname"

    if [ "$1" ]; then
	target=$1
    else
	print_menu
	echo -n "Choose a target platform (Enter for local machin):"
	read target
    fi

    if [ ! -n "$target" ]; then
	set_compile_env_local
	return 0
    fi

    check_target $target
    if [ $? -ne 0 ]
    then
	echo 
	echo  "** $target is not suportted now!"
	print_menu
	return -1
    fi

    set_compile_env $target

    print_env
    echo 
}

function check_target()
{
    local i
    
    if [ $MITHRANDIR_LOCAL eq 1 ]; then
	return 0
    fi

    for i in ${MITHRANDIR_TARGTS[@]}; do
	if [ "$1" = "$i" ]; then
	    return 0
	fi
    done

    return 1
}

unset MITHRANDIR_TARGTS
function add_new_target()
{
    local new_target=$1
    local i
    
    for i in ${MITHRANDIR_TARGTS[@]}; do
	if [ "$new_target" = "$i" ]; then
	    return
	fi
    done
    MITHRANDIR_TARGTS=(${MITHRANDIR_TARGTS[@]} $new_target)
}

function print_menu()
{
    local target=1
    echo "Currently supported targets:"
    for target in ${MITHRANDIR_TARGTS[@]}
    do
	echo "	    $target"
    done

    echo
}

function _mithrandir()
{
    local cur prev opts
    REPLY=()
    cur="${TARGET_WORDS[COMP_CWORD]}"
    prev="${TARGET_WORDS[COMP_CWORD-1]}"

    REPLY=( $(compgen -W "${MITHRANDIR_TARGTS[*]}" -- ${cur}) )
    return 0
}
complete -F _mithrandir mithrandir

function print_env()
{
    local toolchain
    if [ ! -n "$CROSS_COMPILE" ]; then
	toolchain="Use local toolchain"
    else
	toolchain="${CROSS_COMPILE##*/}"
    fi

    echo 
    echo "*******************************************************"
    echo
    echo "Project top directory is:  $MITHRANDIR_TOP"
    echo "Target platform is:        $ARCH"
    echo "Toolchain is:              $toolchain"
    echo
    echo "*******************************************************"
}

function set_compile_env()
{
    T=$(gettop)

    export ARCH=$1

    case $ARCH in
	x86_64) export toolchain=x86_64-unknown-linux-gnu
	    ;;
	arm) export toolchain=arm-unknown-linux-gnueabi
	    ;;
	*) echo "Can't find toolchain for unknown architecture: $ARCH"
	    ;;
    esac

    export CROSS_COMPILE=$T/tools/toolchain/$ARCH/$toolchain/bin/$toolchain-
    export MITHRANDIR_TOP=$T
    export GCC_COLORS='error=01;31:warning=01;35:note=01;36:caret=01;32:locus=01:quote=01'
}

function set_compile_env_local()
{
    export ARCH=$(uname -m)
    unset CROSS_COMPILE
    export MITHRANDIR_TOP=$(gettop)
    export GCC_COLORS='error=01;31:warning=01;35:note=01;36:caret=01;32:locus=01:quote=01'
    print_env
}

function gettop
{
    local TOPFILE=script/setenv.sh
    if [ -n "$TOP" -a -f "$TOP/$TOPFILE" ] ; then
        echo $TOP
    else
        if [ -f $TOPFILE ] ; then
            # The following circumlocution (repeated below as well) ensures
            # that we record the true directory name and not one that is
            # faked up with symlink names.
            PWD= /bin/pwd
        else
            local HERE=$PWD
            T=
            while [ \( ! \( -f $TOPFILE \) \) -a \( $PWD != "/" \) ]; do
                \cd ..
                T=`PWD= /bin/pwd -P`
            done
            \cd $HERE
            if [ -f "$T/$TOPFILE" ]; then
                echo $T
            fi
        fi
    fi
}

add_new_target arm
add_new_target x86_64

mithrandir 
