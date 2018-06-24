while read line || [ -n "$line" ]; do
  if [ -z $n ]; then
    n=${#line}
    for ((i=0; i<${#line}; i++)); do
      if [ ${line:$i:1} == "C" ]; then
        c=$i
        break
      elif [ ${line:$i:1} == "_" ]; then
        c=$i
      fi
    done
    x="|"
  else
    for ((i=$p-1; i<=$p+1; i++)); do
      if [ $i -lt 0 ]; then
        continue
      elif [ $i -ge $n ]; then
        break
      elif [ ${line:$i:1} == "C" ]; then
        c=$i
        break
      elif [ ${line:$i:1} == "_" ]; then
        c=$i
      fi
    done
    if [ $c -lt $p ]; then
      x="/"
    elif [ $c -eq $p ]; then
      x="|"
    else
      x="\\"
    fi
  fi
  line=${line:0:$c}$x${line:$(($c+1))}
  echo "$line"
  p=$c
done <$1
