while read -r line || [ -n "$line" ]; do
  echo "$line"|tr "[:lower:][:upper:]" "[:upper:][:lower:]"
done <$1
