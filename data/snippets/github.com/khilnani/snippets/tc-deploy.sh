#!/bin/sh
# @author nkhilnani


if [ "$1" == "" ]
then

    echo ""
    echo ""
    echo "     USAGE 1: $0 [SERVER_NAME] [TC_NAME] [WAR_DIR] [WAR_NAME]"
    echo ""
    echo "     AWS AMI TC Dir: /usr/share/"
    echo "     Ensure tomcat does not explode Wars"
    echo ""
    
	echo -n "SERVER_NAME: "
	read -e SN
	
	echo -n "TC_NAME: "
	read -e TC
	
	echo -n "WAR_DIR: "
	read -e WRDIR

	echo -n "WAR_NAME: "
	read -e WR
	
	
else
	
	SN="$1"
	TC="$2"
	WRDIR="$3"
	WR="$4"

fi

	today="`eval date +%Y-%m-%d_%H-%M-%S`"

	war="$TC"

	if [ "$WR" != "" ]
	then
		war="$WR"
	fi
	
    echo ""
    echo "     Using file: $WRDIR/$WR.war"
    echo "     Backing up existing WAR to: $war.war.$today.$USER"	
    echo ""

	scp $WRDIR/$war.war $USER@$SN:/tmp/$war.war.$USER

	{ sleep 1
		echo cd /tmp
		echo pwd
		echo ls -la
		echo chmod 755 $war.war.$USER
		sleep 2
		echo sudo -u apache cp /usr/share/$TC/webapps/$war.war /usr/share/$TC/webapps/$war.war.$today.$USER
		sleep 2
		echo sudo -u apache cp $war.war.$USER /usr/share/$TC/webapps/$war.war
		sleep 2
		echo cd /usr/share/$TC/webapps/
		sleep 2
		echo pwd
		echo ls -la $war.war
	} | ssh $USER@$SN -t -T

