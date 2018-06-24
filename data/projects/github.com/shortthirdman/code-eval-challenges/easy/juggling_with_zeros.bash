while read line || [ -n "$line" ]; do
  a=( $line )
  r=0
  for ((i=0; i<${#a[*]}; i+=2)); do
    ((r<<=${#a[$i+1]}))
    if [ ${#a[$i]} -eq 2 ]; then
      ((r+=(1<<${#a[$i+1]})-1))
    fi
  done
  echo "$r"
done <$1
