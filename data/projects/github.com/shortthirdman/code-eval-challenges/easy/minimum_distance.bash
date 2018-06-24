while read line || [ -n "$line" ]; do
  IFS=$' ' a=( $line )
  n=${a[0]}
  a=("${a[@]:1}")
  IFS=$'\n' so=($(sort -n <<<"${a[*]}"))
  u=${so[$(($n/2))]}
  n=0
  for a in ${so[*]}; do
    if [ $a -gt $u ]; then
      ((n+=$a-$u))
    else
      ((n+=$u-$a))
    fi
  done
  echo "$n"
done <$1
