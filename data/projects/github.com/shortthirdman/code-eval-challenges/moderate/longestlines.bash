read n <$1
while read -r line || [ -n "$line" ]; do
  for (( i=1; i<=$n; i++ )); do
    if [ ${#line} -gt ${#a[$i]} ]; then
      a[$i-1]=${a[$i]}
      a[$i]=$line
    fi
  done
done <$1
for (( i=$n; i>0; i-- )); do
  echo "${a[$i]}"
done
