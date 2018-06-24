sed -r ':A;s|([0-9])([0-9])|\1 \2|g;t A' $1|while read line || [ -n "$line" ]; do
  a=( $line )
  for ((i=${#a[*]}-2; i>=0; i-=2)); do
    a[$i]=$((a[$i]*2))
    if [ ${a[$i]} -gt 9 ]; then
      a[$i]=$((a[$i]%10+1))
    fi
  done
  s=0
  for i in ${a[*]}; do
    ((s+=$i))
  done
  if [ $(($s%10)) -eq 0 ]; then
    echo 1
  else
    echo 0
  fi
done
