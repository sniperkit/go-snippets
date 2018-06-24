while read line || [ -n "$line" ]; do
  a=( $line )
  if [ ${#a[*]} -gt 1 ]; then
    echo "${a[${#a[*]}-2]}"
  else
    echo
  fi
done <$1
