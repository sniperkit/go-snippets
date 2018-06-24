p=( ", yeah!" ", this is crazy, I tell ya."  ", can U believe this?" ", eh?" ", aw yea." ", yo." "? No way!" ". Awesome!" )
l=false
c=0
while read -r line || [ -n "$line" ]; do
  for ((i=0; i<${#line}; i++)); do
    if [ "${line:$i:1}" == '.' ] || [ "${line:$i:1}" == '!' ] || [ "${line:$i:1}" == '?' ]; then
      if $l; then
        printf "${p[$c]}"
        l=false
        c=$((($c+1)%8))
      else
        printf "${line:$i:1}"
        l=true
      fi
    else
      printf "${line:$i:1}"
    fi
  done
  echo
done <$1
