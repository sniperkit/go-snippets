#!/bin/sh
#########################################################
# Installer for udev usb rules for orbbec devices
#########################################################

# check if root
if [ $(id -u) -ne 0 ]; then
    # escalate
    echo "\nInstaller requires root, escalating via sudo."
    sudo "$0" "$@"
    exit $?
fi

echo "Linux installer script for Astra SDK\n"

CURRENT_PATH=`pwd`
cd `dirname $0`
# save installer script location
SDK_PATH=`pwd`
# return to original path
cd $CURRENT_PATH

# if not OSX, install UDEV usb rules
if [ "`uname -s`" != "Darwin" ]; then
    echo "Installing rules for orbbec devices into /etc/udev/rules.d/"

    # Install UDEV rules for USB device
    if [ -f ${SDK_PATH}/orbbec-usb.rules ]; then
        cp ${SDK_PATH}/orbbec-usb.rules /etc/udev/rules.d/558-orbbec-usb.rules
    else
        echo "!!!!!!! orbbec-usb.rules file not found. please re-extract Astra SDK"
    fi
fi

echo "Done.\n"

echo "NOTES:"
echo "We suggest adding the following lines to your .bash_profile or .bashrc"
echo "export ASTRA_SDK_INCLUDE=$SDK_PATH/include"
echo "export ASTRA_SDK_LIB=$SDK_PATH/lib\n\n"
