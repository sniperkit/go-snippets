while read line || [ -n "$line" ]; do
  a=( $line )
  n=${a[-1]}
  unset a[${#a[@]}-1]
  unset a[${#a[@]}-1]
  for ((i=0; i<${#a[@]}-1 && i<$n; i++)); do
    for ((j=1; j<${#a[@]}; j++)); do
      if [ ${a[$j-1]} -gt ${a[$j]} ]; then
        x=${a[$j-1]}
        a[$j-1]=${a[$j]}
        a[$j]=$x
      fi
    done
  done
  echo "${a[*]}"
done <$1
