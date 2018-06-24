while read line || [ -n "$line" ]; do
    a=( `echo "$line" | tr ' ' '\n' | sort -g` )
    echo "${a[*]}"
done <$1
