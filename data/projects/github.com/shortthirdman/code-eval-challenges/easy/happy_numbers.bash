while read line || [ -n "$line" ]; do
    b=( $line )
    for ((i=128; i>0; i--)); do
        if [ $line -le 1 ]; then
            break
        fi
        a=$line
        line=
        while [ $a -gt 0 ]; do
            r=$(($a%10))
            ((line+=$r*$r))
            ((a/=10))
        done
        for j in $b; do
            if [ $j -eq $line ]; then
                line=
                break
            fi
        done
        b+=( $line )
    done
    if [ $line -eq 1 ]; then
        echo "1"
    else
        echo "0"
    fi
done <$1
