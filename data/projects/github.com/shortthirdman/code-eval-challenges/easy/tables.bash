for ((i=1; i<13; i++)); do
  for ((j=1; j<13; j++)); do
    printf "%4d" $(($i*$j))
  done
  echo
done
