while read line || [ -n "$line" ]; do
  n=$((${#line}/2))
  r=1
  for ((i=0; i<$n; i++)); do
    x=${line:$i:1}
    y=${line:$(($i+$n)):1}
    if [[ $x == A ]] && [[ $y == B ]]; then
      r=0
      break
    elif [[ $x == B ]] && [[ $y == A ]]; then
      r=0
      break
    elif [[ $x == \* ]] && [[ $y == \* ]]; then
      ((r*=2))
    fi
  done
  echo "$r"
done <$1
