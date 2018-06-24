while read line || [ -n "$line" ]; do
  IFS=$' ' a=( $line )
  IFS=$'\n' b=($(sort -r <<<"${a[*]}"))
  echo ${b[*]}
done <$1
