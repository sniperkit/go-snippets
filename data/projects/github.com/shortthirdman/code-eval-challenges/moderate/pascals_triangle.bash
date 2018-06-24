while read line || [ -n "$line" ]; do
    printf 1
    for ((i=1; i<$line; i++)); do
        printf " 1"
        r=1
        for ((j=1; j<=i; j++)); do
            r=$(( ($r*($i+1-$j))/$j ))
            printf " %d" $r
        done
    done
    echo
done <$1
