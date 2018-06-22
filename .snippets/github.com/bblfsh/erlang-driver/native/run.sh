#!/bin/bash

#This script is used only for development purposes
#Run.sh builds the native part of the driver without anything related with bblfs-sdk
rebar3 release

sed -i 's/-noinput/ /'  _build/default/rel/pre-alpha/bin/pre-alpha
sed -i '/echo "Exec:/ d' _build/default/rel/pre-alpha/bin/pre-alpha
sed -i '/echo "Root:/ d' _build/default/rel/pre-alpha/bin/pre-alpha
sed -i '/echo "$RELEASE_ROOT_DIR"/ d' _build/default/rel/pre-alpha/bin/pre-alpha


exec _build/default/rel/pre-alpha/bin/pre-alpha foreground
