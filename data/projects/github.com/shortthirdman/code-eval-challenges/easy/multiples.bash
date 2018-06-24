while read line || [ -n "$line" ]; do
  a=( ${line//,/ } )
  i=0
  while [ $((${a[1]}*$i)) -lt ${a[0]} ]; do
    ((i++))
  done
  echo "$((${a[1]}*$i))"
done <$1
