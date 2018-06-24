while read line || [ -n "$line" ]; do
  a=( $line )
  n=${a[-1]}
  unset a[${#a[@]}-1]
  unset a[${#a[@]}-1]
  l=1
  for ((k=1; k<=$n; k++)); do
    p=$l
    l=
    for ((i=$p; i<${#a[@]}; i++)); do
      if [ ${a[$i-1]} -gt ${a[$i]} ]; then
        x=${a[$i-1]}
        a[$i-1]=${a[$i]}
        a[$i]=$x
        if [ $i -gt 1 ]; then
          l=$(($i-1))
        else
          l=2
        fi
        break
      fi
    done
    if [ -z $l ]; then
      break
    fi
  done
  echo "${a[*]}"
done <$1
