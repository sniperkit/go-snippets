while read line || [ -n "$line" ]; do
  l=
  for i in $line; do
    if [ ${#i} -gt ${#l} ]; then
      l=$i
    fi
  done
  for ((i=0; i<${#l}; i++)); do
    if [ $i -gt 0 ]; then
      printf " "
    fi
    for ((j=0; j<$i; j++)); do
      printf "*"
    done
    printf ${l:$i:1}
  done
  echo
done <$1
