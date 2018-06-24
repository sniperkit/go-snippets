while read -r line || [ -n "$line" ]; do
    echo "${line,,}"
done <$1
