while read line || [ -n "$line" ]; do
  a=( ${line//|/ } )
  c=0
  for ((i=0; i<${#a[*]}/2-2; i++)); do
    for ((j=$i+1; j<${#a[*]}/2-1; j++)); do
      ((dx=${a[$i*2]}-${a[$j*2]}))
      ((dy=${a[$i*2+1]}-${a[$j*2+1]}))
      for ((k=0; k<${#a[*]}/2; k++)); do
        if [[ $k != $i ]] && [[ $k != $j ]] && [[ $(($dx*(${a[$i*2+1]}-${a[$k*2+1]}))) = $(((${a[$i*2]}-${a[$k*2]})*$dy)) ]]; then
          if [[ $k -gt $j ]]; then
            ((c++))
          fi
          break
        fi
      done
    done
  done
  echo "$c"
done <$1
