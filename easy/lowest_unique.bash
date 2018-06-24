while read line || [ -n "$line" ]; do
  a=( $line )
  u=0
  v=0
  for i in ${a[*]}; do
    if [ $((u&(1<<$i))) -gt 0 ]; then
      ((v|=(1<<$i)))
    else
      ((u|=(1<<$i)))
    fi
  done
  i=1
  while [ $u -gt 0 ]; do
    if [ $((u&(1<<$i))) -gt 0 ] && [ $((v&(1<<$i))) -eq 0 ]; then
      break
    fi
    ((u&=(1022-(1<<$i))))
    ((i++))
  done
  if [ $u -eq 0 ]; then
    echo "0"
  else
    r=0
    while [ ${a[$r]} -ne $i ]; do
      ((r++))
    done
    ((r++))
    echo "$r"
  fi
done <$1
