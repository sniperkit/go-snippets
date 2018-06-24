sed -e "/^$/d; s/^{/0/; s/{\"id\": [^,}]\+},\? \?//g; s/{\"id\": \([^,]\+\),[^}]\+}/\1 /g; s/[^ 0123456789]//g" $1 | while read line || [ -n "$line" ]; do
  a=( $line )
  s=0
  for i in ${a[*]}; do
    ((s+=$i))
  done
  echo "$s"
done
