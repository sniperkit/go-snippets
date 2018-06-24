while read line || [ -n "$line" ]; do
  a=( $line )
  if [ ${#a[*]} -gt 1 ]; then
    s0=${a[0]}
    s1=${a[1]}
    c=0
    o=1
    r=0
    v=0
    for (( i=0; i<${#s1}; i++ )); do
      if [[ "${s1:$i:1}" =~ [a-z] ]]; then
        ((v=v*10+${s0:$c:1}))
        ((c++))
      elif [ "${s1:$i:1}" = "+" ]; then
        ((r+=v*o))
        o=1
        v=0
      else
        ((r+=v*o))
        o=-1
        v=0
      fi
    done
    echo "$((r+v*o))"
  fi
done <$1
