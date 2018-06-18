#!/bin/sh

if [ "$1" == "" ]
then
    echo ""
    echo "     USAGE 1: $0 SERVER_NAME TC_NAME"
    echo "     USAGE 2: $0 SERVER_NAME TC_NAME WAR_NAME (without '.war'"
    echo ""
    echo "     AWS AMI TC Dir: /usr/share/"
    echo ""
    
	echo -n "SERVER_NAME: "
	read -e SN

	echo -n "TC_NAME: "
	read -e TC

	echo -n "WAR_NAME (Without .war): (ENTER to use TC_NAME)"
	read -e WR
	
else
	
	SN="$1"
	TC="$2"
	WR="$3"

fi

	war="$TC"

	if [ "$WR" != "" ]
	then
		war="$WR"
	fi

	{ sleep 1
		echo sudo -u apache /usr/share/$TC/bin/shutdown.sh
		sleep 2
		echo cd /usr/share/$TC/webapps
		sleep 2
		echo ls -la
		sleep 2
		echo sudo -u apache rm -rf $war
		sleep 5
		echo ls -la *.war
		sleep 2
		echo cd /usr/share/$TC/work
		sleep 2
		echo pwd
		sleep 2
		echo sudo -u apache rm -rf Catalina
		sleep 2
		echo ls -la
		sleep 2
		echo sudo -u apache /usr/share/$TC/bin/startup.sh
		sleep 2
		echo tail -f /usr/share/$TC/logs/catalina.out
	} | ssh $USER@$SN -t -T

