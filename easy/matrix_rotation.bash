while read line || [ -n "$line" ]; do
  a=( $line )
  l=1
  while [ $(($l*$l)) -lt ${#a[*]} ]; do
    ((l++))
  done
  for ((i=0; i<l; i++)); do
    for ((j=l-1; j>=0; j--)); do
      printf ${a[j*l+i]}
      if [ $i -lt $(($l-1)) ] || [ $j -gt 0 ]; then
        printf " "
      else
        echo
      fi
    done
  done
done <$1
