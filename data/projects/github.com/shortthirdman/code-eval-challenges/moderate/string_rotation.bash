tr ", " " ," <$1 | while read line || [ -n "$line" ]; do
  a=( $line )
  if [ ${#a[0]} -ne ${#a[1]} ]; then
    echo False
    continue
  fi
  f=false
  for ((i=0; i<${#a[0]}; i++)); do
    if [ ${a[0]} == ${a[1]:$i:$((${#a[0]}-$i))}${a[1]:0:$i} ]; then
      f=true
      break
    fi
  done
  if [ $f ]; then
    echo True
  else
    echo False
  fi
done
