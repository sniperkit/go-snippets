tr -cd "[:digit:] \n" <$1 | while read line || [ -n "$line" ]; do
  a=( $line )
  if [ $((a[0] + a[1] + a[2])) -eq 0 ]; then
    echo "0"
  else
    echo $(( (a[0] * 3 + a[1] * 4 + a[2] * 5) * a[3] / (a[0] + a[1] + a[2]) ))
  fi
done
