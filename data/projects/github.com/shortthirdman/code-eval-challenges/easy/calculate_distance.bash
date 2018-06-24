tr -d "[()," <$1 | while read line || [ -n "$line" ]; do
  a=( $line )
  x=$((${a[0]}-${a[2]}))
  y=$((${a[1]}-${a[3]}))
  i=0
  while [ $(($x*$x+$y*$y)) -ge $(($i*$i)) ]; do
    ((i++))
  done
  echo "$(($i-1))"
done
