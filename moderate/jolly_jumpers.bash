while read line || [ -n "$line" ]; do
  a=( $line )
  j=true
  if [ $a -gt 1 ]; then
    declare -a u
    for ((i=2; i<${#a[*]}; i++)); do
      if [ ${a[$i]} -lt ${a[$i-1]} ]; then
        x=$((${a[$i-1]}-${a[$i]}))
      else
        x=$((${a[$i]}-${a[$i-1]}))
      fi
      if [ $x -ge $a ] || [ $x -eq 0 ] || [[ "${u[$x]}" == "true" ]]; then
        j=false
        break
      fi
      u[$x]="true"
    done
    unset u
  fi
  if [ "$j" = true ]; then
    echo "Jolly"
  else
    echo "Not jolly"
  fi
done <$1
