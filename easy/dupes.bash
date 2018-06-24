while read line || [ -n "$line" ]; do
    a=( ${line//,/ } )
    printf $a
    for ((i=1; i<${#a[*]}; i++)); do
        if [ ${a[$i-1]} -ne ${a[$i]} ]; then
            printf ",%d" ${a[$i]}
        fi
    done
    echo
done <$1
