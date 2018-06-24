while read line || [ -n "$line" ]; do
  if [ -n "$line" ]; then
    a=( $line )
    i=$((${#a[@]}-1))
    printf ${a[$i]}
    for ((i--; i>=0; i--)); do
        printf " ${a[$i]}"
    done
    echo
  fi
done <$1
