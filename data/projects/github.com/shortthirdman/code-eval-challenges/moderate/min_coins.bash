b=( 0 1 2 1 2 )
while read line || [ -n "$line" ]; do
  echo $(( $line/5 + ${b[$line%5]} ))
done <$1
