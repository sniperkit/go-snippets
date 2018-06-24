sed -e "s/, / /g" $1| while read line || [ -n "$line" ]; do
  a=( $line )
  l=$(( ${#a[@]} - 1 ))
  m=$(( ${#line} - ${#a[$l]} - 1 ))
  echo "$line" |cut -c "-$m" |tr -d "${a[$l]}" |sed -e "s/^[[:space:]]*//" -e "s/[[:space:]]*$//"
done
