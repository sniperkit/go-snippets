while read line || [ -n "$line" ]; do
  if [[ $line == *" C1970"* ]]; then
    a=( $line )
    if [ $a == "C" ]; then
      chep=( ${a[2]:2} ${a[3]:2} ${a[4]//[![:digit:]]} )
    else
      ((fit=${a[2]:2}*${chep[0]}+${a[3]:2}*${chep[1]}+${a[4]:2}*${chep[2]}))
      for (( i=1; i<=6; i++ )); do
        if [ -z "${b[$i]}" ] || [ $fit -gt ${b[$i]} ]; then
          b[$i-1]=${b[$i]}
          c[$i-1]=${c[$i]}
          b[$i]=$fit
          c[$i]=${a[1]:1}
        fi
      done
    fi
  fi
done <$1
sum=${c[1]}
for (( i=2; i<=6; i++ )); do
  ((sum+=${c[i]}))
done
echo "$sum"
