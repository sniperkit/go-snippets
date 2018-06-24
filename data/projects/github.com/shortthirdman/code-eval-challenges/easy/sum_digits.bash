while read line || [ -n "$line" ]; do
  res=
  for ((i=0; i<${#line}; i++)); do
    res=$(($res+${line:$i:1}))
  done
  echo "$res"
done <$1
