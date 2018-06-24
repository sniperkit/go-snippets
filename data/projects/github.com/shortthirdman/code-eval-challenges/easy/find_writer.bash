tr "|" "\n" <$1 | while read line || [ -n "$line" ]; do
  if [ -n "$line" ]; then
    a="$line"
    read line
    b=( $line )
    for i in ${b[*]}; do
      printf "${a:$(( $i-1 )):1}"
    done
    echo
  fi
done
