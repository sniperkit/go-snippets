while read line || [ -n "$line" ]; do
    echo "$((($line + 1) % 2))"
done <$1
