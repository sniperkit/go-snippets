while read line || [ -n "$line" ]; do
  a=( ${line//[x|]/ } )
  r=0
  while [[ ${a[1]} != ${a[3]} ]]; do
    t=${a[0]}
    ((r+=$t))
    ((a[0]=${a[1]}-1))
    a[1]=$t
    t=${a[2]}
    ((a[2]=${a[0]}+1-${a[3]}))
    a[3]=$t
  done
  echo $(($r+${a[2]}))
done <$1
