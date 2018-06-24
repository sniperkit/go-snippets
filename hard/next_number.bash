dsig() {
  u=0
  for ((i=0; i<${#1}; i++)); do
    if [[ ${1:$i:1} > 0 ]]; then
      ((u+=1<<(3*${1:$i:1})))
    fi
  done
}

while read line || [ -n "$line" ]; do
  ((b=line+9))
  dsig $line
  d=$u
  while dsig $b; [[ $d != $u ]]; do
    ((b+=9))
  done
  echo "$b"
done <$1
