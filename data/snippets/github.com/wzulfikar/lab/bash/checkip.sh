# check public IP
ip() {
	LASTIP="$HOME/lastipaddr"
	GETIPADDR="checkip.amazonaws.com"
	CURRENTIP=$(curl $GETIPADDR 2>/dev/null)
	if [ -f $LASTIP ]
	then
	   # compare last word in file $LASTIP with $CURRENTIP
	   if [ `awk 'END {print $NF}' $LASTIP` = $CURRENTIP ]
	   then
	       echo "$CURRENTIP (Unchanged)"
	   else
	       echo "$(date "+%F %T %Z") - $CURRENTIP" >> $LASTIP
	       echo "$CURRENTIP (Changed)"
	    fi
	else
	    echo $CURRENTIP >> $LASTIP
	fi
}

# check last 3 IPs
iplast() {
	LASTIP="$HOME/lastipaddr"
	tail -n 3 $LASTIP
}