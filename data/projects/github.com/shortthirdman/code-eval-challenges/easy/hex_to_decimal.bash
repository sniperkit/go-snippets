while read line || [ -n "$line" ]; do
    echo "$((0X$line))"
done <$1
