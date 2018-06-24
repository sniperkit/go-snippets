while read line || [ -n "$line" ]; do
  a=( ${line/;/ } )
  declare -a r
  h=$(( ${#a[*]}/2+1 ))
  for ((i=0; i<$h-1; i++)); do
    r[${a[$h+$i]}]=${a[$i]}
  done
  for ((i=1; i<=$h; i++)); do
    if [ -z ${r[$i]} ]; then
      r[$i]=${a[$(( $h-1 ))]}
    fi
  done
  echo "${r[*]}"
  unset a r
done <$1
