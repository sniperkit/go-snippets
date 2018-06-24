ord() {
  printf -v v '%d' "'$1"
}

while read line || [ -n "$line" ]; do
  declare -a ch
  for ((c=0; c<${#line}; c++)); do
    ord ${line:$c:1}
    ((ch[$v]++))
  done
  for ((c=0; c<${#line}; c++)); do
    ord ${line:$c:1}
    if [ ${ch[$v]} -eq 1 ]; then
      printf ${line:$c:1}
      break
    fi
  done
  echo
  unset ch
done <$1
