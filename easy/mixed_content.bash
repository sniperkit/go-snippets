while read -r line || [ -n "$line" ]; do
  a=( ${line//,/ } )
  declare -a n s
  for i in ${a[*]}; do
    if [[ $i =~ ^[0-9]+$ ]]; then
      n=( ${n[*]} $i )
    else
      s=( ${s[*]} $i )
    fi
  done
  if [ ${#s[*]} -gt 0 ]; then
    printf $s
    for ((i=1; i<${#s[*]}; i++)); do
      printf ","${s[$i]}
    done
    if [ ${#n[*]} -gt 0 ]; then
      printf "|"
    fi
  fi
  if [ ${#n[*]} -gt 0 ]; then
    printf $n
    for ((i=1; i<${#n[*]}; i++)); do
      printf ","${n[$i]}
    done
  fi
  echo
  unset n s
done <$1
