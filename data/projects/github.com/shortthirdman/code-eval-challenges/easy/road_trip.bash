tr -cd "[:digit:] \n" <$1 | while read line || [ -n "$line" ]; do
    a=( `echo "$line" | tr ' ' '\n' | sort -n -S 1M` )
    printf ${a[0]}
    for ((c=1; c<${#a[*]}; c++)); do
        printf ",%d" $((${a[$c]}-${a[$c-1]}))
    done
    echo
done
