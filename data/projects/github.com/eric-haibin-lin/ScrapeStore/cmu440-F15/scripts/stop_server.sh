while read line
do
    kill -9 $line 2> /dev/null
done <server_log.txt 2> /dev/null
rm server_log.txt 2> /dev/null
