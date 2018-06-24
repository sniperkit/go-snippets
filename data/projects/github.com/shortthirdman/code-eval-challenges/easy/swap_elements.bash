tr ",-" " " <$1 | while read line || [ -n "$line" ]; do
  a=( $line )
  c=-1
  d=-1
  ix=0
  for i in ${a[*]}; do
    if [ $c -ne -1 ]; then
      if [ $d -eq -1 ]; then
        d=$i
      else
        t=${a[$d]}
        a[$d]=${a[$i]}
        a[$i]=$t
        d=-1
      fi
    elif [ $i == ":" ]; then
      c=$ix
    else
      ((ix++))
    fi
  done
  printf $a
  for ((i=1; i<ix; i++)); do
    printf " "${a[$i]}
  done
  echo
done
