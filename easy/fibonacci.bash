while read line || [ -n "$line" ]; do
    if [ $line -eq 0 ] || [ $line -eq 1 ]; then
        echo "$line"
    else
        b=( 1 1 0 )
        c=1
        while [ $line -gt $c ]; do
            ((c++))
            b[0]=$((${b[1]}+${b[2]}))
            b[2]=${b[1]}
            b[1]=${b[0]}
        done
        echo "${b[0]}"
    fi
done <$1
