while read line || [ -n "$line" ]; do
    b=-1
    c=1
    for a in $line; do
        if [ $a -eq $b ]; then
            ((c++))
        else
            if [ $b -ge 0 ]; then
                printf "%d %d " $c $b
                c=1
            fi
            b=$a
        fi
    done
    printf "%d %d\n" $c $b
done < $1
