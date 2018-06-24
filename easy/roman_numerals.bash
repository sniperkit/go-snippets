ronum=( 1000 900 500 400 100 90 50 40 10 9 5 4 1 )
rostr=( M CM D CD C XC L XL X IX V IV I )
while read line || [ -n "$line" ]; do
  i=0
  while [ $line -gt 0 ]; do
    if [ $line -ge ${ronum[$i]} ]; then
      printf ${rostr[$i]}
      ((line-=${ronum[$i]}))
    else
      ((i++))
    fi
  done
  echo
done <$1
