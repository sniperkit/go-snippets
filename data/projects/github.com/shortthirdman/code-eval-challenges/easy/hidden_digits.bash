tr -cd "[:digit:]a-j\n" <$1 | while read line || [ -n "$line" ]; do
  if [ ${#line} -eq 0 ]; then
    echo "NONE"
  else
    echo "$line" | tr "a-j" "0-9"
  fi
done
