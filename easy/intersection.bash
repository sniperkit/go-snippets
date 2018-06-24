tr ",;" " \n" <$1 | while read line || [ -n "$line" ]; do
  a=( $line )
  read line
  b=( $line )
  j=0
  r=
  for ((i=0; i<${#a[*]}; i++)); do
    while [ $j -lt ${#b[*]}  ] && [ ${a[$i]} -gt ${b[$j]} ]; do
      ((j++))
    done
    if [ $j -lt ${#b[*]}  ] && [ ${a[$i]} -eq ${b[$j]} ]; then
      if [ -z $r ]; then
        r=${a[$i]}
      else
        r="$r,${a[$i]}"
      fi
      ((j++))
    fi
  done
  echo "$r"
done
