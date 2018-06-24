while read line || [ -n "$line" ]; do
  declare -a a
  n=$line
  r=0
  self=1
  for ((i=0; i<${#n}; i++)); do
    a[$i]=0
  done
  while [ $line -ne 0 ]; do
    i=$(($line%10))
    if [ $i -le ${#n} ]; then
      ((a[$i]++))
    else
      self=0
      break
    fi
    ((line/=10))
  done
  if [ $self -eq 1 ]; then
    for i in ${a[*]}; do
      r=$(($r*10+$i))
    done
    if [ $n -ne $r ]; then
      self=0
    fi
  fi
  echo "$self"
  unset a
done <$1
