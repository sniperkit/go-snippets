#!/bin/bash

if [ "${AUTHORIZED_KEYS}" == "" -a "${ROOT_PASS}" == "" ]; then
  ROOT_PASS=$(pwgen 10 --capitalize --numerals --ambiguous --symbols | sed 's/"//')
elif echo "${ROOT_PASS}" | grep -Ew ".{1,9}" > /dev/null; then
  echo ""
  echo "WARNING: Your ROOT_PASS is very weak!"
  echo "WARNING: Such password can crack in a few hours!"
  echo ""
  echo "ROOT_PASS Recommended:"
  echo "  * Minimum 10 characters in length"
  echo "  * Contains 3/4 of the following items:"
  echo "    - Uppercase Letters ... [A-Z]"
  echo "    - Lowercase Letters ... [a-z]"
  echo "    - Numbers ... [0-9]"
  echo '    - Symbols ... [!#$%&()*+,-./:;<=>?@[\]^_`{|}]'
  echo ""
  echo "INFO: We have rewritten your ROOT_PASS with a strong password."

  ROOT_PASS=$(pwgen 10 --capitalize --numerals --ambiguous --symbols | sed 's/"//')
fi


echo ""
echo "          _[]___                        "
echo "       __|_o_o_o\__                     "
echo '       \""""""""""/                     '
echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
echo ""
echo "SSH Connection Information:"
echo ""
echo "    ssh -p <port> root@<host>"
echo ""

if [ "${ROOT_PASS}" != "" ]; then
  echo "root:$ROOT_PASS" | chpasswd > /dev/null
  sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config
  echo "Then enter the root password: ${ROOT_PASS}"
  echo ""
fi

echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
echo ""
