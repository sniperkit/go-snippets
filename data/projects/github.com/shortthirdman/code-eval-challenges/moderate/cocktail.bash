while read line || [ -n "$line" ]; do
  a=( $line )
  n=${a[-1]}
  unset a[${#a[@]}-1]
  unset a[${#a[@]}-1]
  ((n=($n<${#a[*]}/2)?$n:(${#a[*]}/2)))
  for ((k=1; k<=$n; k++)); do
    for ((i=$k; i<=${#a[@]}-$k; i++)); do
      if [ ${a[$i-1]} -gt ${a[$i]} ]; then
        x=${a[$i-1]}
        a[$i-1]=${a[$i]}
        a[$i]=$x
      fi
    done
    for ((i=${#a[@]}-$k-1; i>=$k; i--)); do
      if [ ${a[$i-1]} -gt ${a[$i]} ]; then
        x=${a[$i-1]}
        a[$i-1]=${a[$i]}
        a[$i]=$x
      fi
    done
  done
  echo "${a[*]}"
done <$1
