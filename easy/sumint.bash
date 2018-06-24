while read line || [ -n "$line" ]; do
  ((s+=$line))
done <$1
echo "$s"
