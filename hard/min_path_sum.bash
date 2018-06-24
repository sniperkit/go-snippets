while read line || [ -n "$line" ]; do
  if [ -z $n ]; then
    n=$line
    i=$n
    declare -a b
  else
    s=( ${line//,/ } )
    if [ $i -eq $n ]; then
      b[0]=${s[0]}
      for ((j=1; j<$n; j++)); do
        b[$j]=$((${b[$j-1]}+${s[$j]}))
      done
    else
      b[0]=$((${b[0]}+${s[0]}))
      for ((j=1; j<$n; j++)); do
        if [ ${b[$j]} -lt ${b[$j-1]} ]; then
          b[$j]=$((${b[$j]}+${s[$j]}))
        else
          b[$j]=$((${b[$j-1]}+${s[$j]}))
        fi
      done
    fi
    ((i--))
    if [ $i -eq 0 ]; then
      echo "${b[$n-1]}"
      n=
      unset b
    fi
  fi
done <$1
