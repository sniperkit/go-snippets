while read line || [ -n "$line" ]; do
  a=( $line )
  n=0
  while [ ${a[$n]} != "|" ]; do
    ((n++))
  done
  for ((o=$n+1; o<${#a[*]}; o+=$n+1)); do
    for ((i=0; i<$n; i++)); do
      if [ ${a[$i]} -lt ${a[$o+$i]} ]; then
        a[$i]=${a[$o+$i]}
      fi
    done
  done
  echo "${a[*]:0:$n}"
done <$1
