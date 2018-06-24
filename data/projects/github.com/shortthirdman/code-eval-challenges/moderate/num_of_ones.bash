while read line || [ -n "$line" ]; do
  i=0
  while [ $line -gt 0 ]; do
    ((line&=$line - 1))
    ((i++))
  done
  echo $i
done <$1
