while read line || [ -n "$line" ]; do
    a=( $line )
    for ((i=0; i<${#a[*]}/2; i++)); do
        if [ $i -gt 0 ]; then
            printf " "
        fi
        printf $((${a[i]}*${a[${#a[*]}/2+i+1]}))
    done
    echo
done <$1
