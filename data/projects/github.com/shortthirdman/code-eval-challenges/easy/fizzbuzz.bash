while read f b n || [ -n "$n" ]; do
  for ((i=1; i<=$n; i++)); do
    if [ $i -gt 1 ]; then
      printf " "
    fi
    if [ $(($i % $f)) -gt 0 ] && [ $(($i % $b)) -gt 0 ]; then
      printf $i
    else
      if [ $(($i % $f)) -eq 0 ]; then
        printf F
      fi
      if [ $(($i % $b)) -eq 0 ]; then
        printf B
      fi
    fi
  done
  echo
done <$1
