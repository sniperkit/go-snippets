#!/bin/sh
eq () {
	NC='\033[0m' # no color
	RED='\033[0;31m'
	GREEN='\033[0;32m'
	ORANGE='\033[0;33m'

	if [ "$1" = "$2" ]
	then
		echo "${GREEN}[✔ IT EQUALS!]${NC}"
		echo "▍$1"
	else
		echo "${RED}[✘ NOT EQUALS!]${NC}"
		echo "▍$1\n${NC}▍$2"
	fi
}
