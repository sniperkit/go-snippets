sed -r ':A;s|([0-9])([0-9])|\1 \2|g;t A' $1|while read line || [ -n "$line" ]; do
  a=( $line )
  for ((i=${#a[*]}-2; i>=0; i-=2)); do
    a[$i]=$((a[$i]*2))
  done
  s=0
  for i in ${a[*]}; do
    ((s+=$i))
  done
  if [ $(($s%10)) -eq 0 ]; then
    echo "Real"
  else
    echo "Fake"
  fi
done
