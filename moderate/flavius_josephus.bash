while read line || [ -n "$line" ]; do
  a=( ${line//,/ } )
  p=0
  declare -a b c
  for ((i=0; i<${a[0]}; i++)); do
    b[$i]=$i
  done
  for ((i=0; i<${a[0]}; i++)); do
    p=$((($p+${a[1]}-1)%${#b[*]}))
    c[$i]=${b[$p]}
    unset b[$p]
    b=(${b[*]})
  done
  echo ${c[*]}
  unset b c
done <$1
