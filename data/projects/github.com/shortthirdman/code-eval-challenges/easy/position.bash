while read line || [ -n "$line" ]; do
  a=( ${line//,/ } )
  if [ $(( !!(${a[0]}&2**(${a[1]}-1)) )) -eq $(( !!(${a[0]}&2**(${a[2]}-1)) )) ]; then
    echo "true"
  else
    echo "false"
  fi
done <$1
