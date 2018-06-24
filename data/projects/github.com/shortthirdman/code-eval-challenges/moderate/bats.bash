while read line || [ -n "$line" ]; do
  a=( $line )
  c=0
  t=0
  tx=$((6-${a[1]}))
  i=6
  while [ $i -le $((${a[0]}-6)) ]; do
    if [ $i -gt $(($tx-${a[1]})) ]; then
      i=$tx
      if [ $t -eq ${a[2]} ]; then
        tx=$((${a[0]}-6+${a[1]}))
      else
        tx=${a[$t+3]}
        ((t++))
      fi
    else
      ((c++))
    fi
    ((i+=${a[1]}))
  done
  echo "$c"
done <$1
