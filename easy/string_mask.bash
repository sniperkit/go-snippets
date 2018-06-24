while read line || [ -n "$line" ]; do
  a=( $line )
  for ((i=0; i<${#a[0]}; i++)); do
    p=${a[0]:$i:1}
    if [ ${a[1]:$i:1} -eq 1 ]; then
      printf ${p^^}
    else
      printf ${p}
    fi
  done
  echo
done <$1
