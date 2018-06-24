while read line || [ -n "$line" ]; do
  if [ -n "$line" ]; then
    c=0
    p=0
    while : ; do
      ((p+=$c))
      c=`expr index "${line:(-1)}" "${line:$p:(-1)}"`
      if [ $c -eq 0 ]; then
        break
      fi
    done
    echo "$(( $p-1 ))"
  fi
done <$1
