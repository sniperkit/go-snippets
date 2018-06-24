while read line || [ -n "$line" ]; do
  p=1
  for ((i=1; i<${#line}; i++)); do
    if [ ${line:$i:1} != ${line:$(( $i-$p )):1} ]; then
      p=$(( $i+1 ))
    fi
  done
  echo "$p"
done <$1
