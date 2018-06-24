while read line || [ -n "$line" ]; do
  a=( $line )
  m=${a[-1]}
  ((n=${#a[*]}-2))
  a=( ${a[*]:0:$n} )
  while [ ${#a[*]} -gt 1 ]; do
    ((n=$m%${#a[*]}-1))
    if [ $n -ge 0 ]; then
      a=( ${a[*]:0:$n} ${a[*]:$n+1} )
    else
      ((n=${#a[*]}-1))
      a=( ${a[*]:0:$n} )
    fi
  done
  echo ${a[0]}
done <$1
