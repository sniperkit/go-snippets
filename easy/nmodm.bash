while read line || [ -n "$line" ]; do
    a=( ${line/,/ } )
    echo "$(($a-($a/${a[1]})*${a[1]}))"
done <$1
