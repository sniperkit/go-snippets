while read line || [ -n "$line" ]; do
    l=
    for i in $line; do
        if [ ${#i} -gt ${#l} ]; then
            l=$i
        fi
    done
    echo "$l"
done <$1
