while read line || [ -n "$line" ]; do
    a=( ${line//,/ } )
    l=0
    m=$a
    for i in ${a[*]}; do
        if [ $(($i+$l)) -gt $m ]; then
            m=$(($i+$l))
        fi
        if [ $(($i+$l)) -gt 0 ]; then
            l=$(($i+$l))
        else
            l=0
        fi
    done
    echo $m
done <$1
