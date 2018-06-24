while read line || [ -n "$line" ]; do
  a=( $line )
  b=( ${a//;/ } )
  c=0
  tw=0
  tn=0
  te=$((${b[1]}-1))
  ts=$((${b[0]}-1))
  i=0
  j=1
  printf ${b[2]}
  while [ $j -le $te ]; do
    while [ $j -le $te ]; do
      printf " ${a[$i*${b[1]}+$j]}"
      ((j++))
    done
    ((j--))
    ((i++))
    ((tn++))
    if [ $i -gt $ts ]; then
      break
    fi
    while [ $i -le $ts ]; do
      printf " ${a[$i*${b[1]}+$j]}"
      ((i++))
    done
    ((i--))
    ((j--))
    ((te--))
    if [ $j -lt $tw ]; then
      break
    fi
    while [ $j -ge $tw ]; do
      printf " ${a[$i*${b[1]}+$j]}"
      ((j--))
    done
    ((j++))
    ((i--))
    ((ts--))
    if [ $i -lt $tn ]; then
      break
    fi
    while [ $i -ge $tn ]; do
      printf " ${a[$i*${b[1]}+$j]}"
      ((i--))
    done
    ((j++))
    ((i++))
    ((tw++))
  done
  echo
done <$1
