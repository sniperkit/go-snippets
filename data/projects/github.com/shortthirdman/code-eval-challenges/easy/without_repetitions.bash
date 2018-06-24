sed -e "s/\(.\)\1*/\1/g" $1 | while read -r line || [ -n "$line" ]; do
    echo "$line"
done
