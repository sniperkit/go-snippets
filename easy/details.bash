sed -e "s/[^,]*\(X\.*Y\)[^,]*/\1/g" $1 | while read line || [ -n "$line" ]; do
  a=( ${line//,/ } )
  r=10
  for i in ${a[*]}; do
    if [ ${#i} -lt $r ]; then
      r=${#i}
    fi
  done
  echo "$(($r-2))"
done
