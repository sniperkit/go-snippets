while read line || [ -n "$line" ]; do
  e=0
  s=0
  for ((i=$line; i>0; i/=10)); do
    ((e++))
  done
  for ((i=$line; i>0; i/=10)); do
    ((s+=($i%10)**e))
  done
  if [ $s -eq $line ]; then
    echo "True"
  else
    echo "False"
  fi
done <$1
