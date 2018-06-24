push() {
  stack+=( "$1" )
}

pop() {
  if [ ${#stack[*]} -gt 0 ]; then
    r=${stack[-1]}
    unset stack[${#stack[*]}-1]
  else
    r=
  fi
}

while read line || [ -n "$line" ]; do
  a=( $line )
  declare -a stack
  for i in ${a[*]}; do
    push $i
  done
  first=true
  while [ ${#stack[*]} -gt 0 ]; do
    if $first; then
      first=false
    else
      printf " "
    fi
    pop
    printf '%s' "$r"
    pop
  done
  echo
  unset stack
done <$1
