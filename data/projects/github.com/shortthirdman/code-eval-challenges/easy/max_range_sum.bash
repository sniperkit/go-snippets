while read line || [ -n "$line" ]; do
  a=( ${line//;/ } )
  m=0
  cm=0
  for ((i=1; i<${#a[*]}; i++)); do
    if [ $i -le $a ]; then
      ((cm+=${a[$i]}))
    else
      ((cm+=${a[$i]}-(${a[$i-$a]})))
    fi
    if [ $i -ge $a ] && [ $cm -gt $m ]; then
      m=$cm
    fi
  done
  echo "$m"
done <$1
