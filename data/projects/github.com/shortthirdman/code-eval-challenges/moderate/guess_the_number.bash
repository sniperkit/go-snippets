while read line || [ -n "$line" ]; do
  a=( $line )
  lo=0
  hi=-1
  for i in $line; do
    if [[ $hi -lt 0 ]]; then
      hi=$i
    elif [[ $i == "Lower" ]]; then
      hi=$(( ($lo+$hi+1)/2-1 ))
    elif [[ $i == "Higher" ]]; then
      lo=$(( ($lo+$hi+1)/2+1 ))
    fi
  done
  echo $(( ($lo+$hi+1)/2 ))
done <$1
