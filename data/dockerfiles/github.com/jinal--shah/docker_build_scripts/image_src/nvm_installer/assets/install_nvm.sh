#!/bin/bash
# vim: et smartindent sr sw=4 ts=4:
#
# install_nvm.sh
#
# - install nvm under /opt/nvm/.nvm
# - that dir should be mounted from host

WORK_DIR=/opt/nvm
INSTALL_DIR=$WORK_DIR/.nvm
if [[ -z "$SRC_URL" ]]; then
    echo "$0 ERROR: \$SRC_URL must be set in env." 2>&1
    exit 1
fi

if [[ ! -r $WORK_DIR ]]; then
    echo "$0 ERROR: $WORK_DIR must be readable dir for installation." 2>&1
    exit 1
fi

git clone $SRC_URL $INSTALL_DIR

if [[ ! -d $INSTALL_DIR ]]; then
    echo "$0 ERROR: unable to successfully clone repo." 2>&1
    exit 1
fi

cd $INSTALL_DIR

LATEST_TAG=$(git describe --abbrev=0 --tags)

if [[ -z "$LATEST_TAG" ]]; then
    echo "$0 ERROR: couldn't find latest release tag" 2>&1
    exit 1
fi

if ! git checkout $LATEST_TAG
then
    echo "$0 ERROR: couldn't check out latest tag $LATEST_TAG" 2>&1
    exit 1
fi

rm -rf $INSTALL_DIR/.git

exit 0

