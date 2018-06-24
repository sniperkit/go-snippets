while read line || [ -n "$line" ]; do
  echo "obase=2;$line"|bc
done <$1
