while read line || [ -n "$line" ]; do
  line=${line/;/,}
  a=( ${line//,/ } )
  n=false
  j=$((${#a[*]}-2))
  for ((i=0; i<$j && 2*${a[i]}<${a[-1]}; i++)); do
    for ((; j>$i && ${a[$i]}+${a[j]}>=${a[-1]}; j--)); do
      if [ $((${a[$i]}+${a[$j]})) -eq ${a[-1]} ]; then
        if $n; then
          printf ';'
        else
          n=true
        fi
        printf "%d,%d" ${a[$i]} ${a[$j]}
      fi
    done
  done
  if $n; then
    echo
  else
    echo "NULL"
  fi
done <$1
