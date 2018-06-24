within() {
  if [ ${a[$2]} -ge ${a[$1]} ] && [ ${a[$2]} -le ${a[$1+2]} ] && [ ${a[$3]} -ge ${a[$1+3]} ] && [ ${a[$3]} -le ${a[$1+1]} ]; then
    printf -v v '%d' 1
  else
    printf -v v '%d' 0
  fi
}

while read line || [ -n "$line" ]; do
  a=( ${line//,/ } )
  within 0 4 5
  if (($v)); then
    echo "True"
    continue
  fi
  within 0 6 7
  if (($v)); then
    echo "True"
    continue
  fi
  within 0 4 7
  if (($v)); then
    echo "True"
    continue
  fi
  within 0 6 5
  if (($v)); then
    echo "True"
    continue
  fi
  within 4 0 1
  if (($v)); then
    echo "True"
    continue
  fi
  within 4 2 3
  if (($v)); then
    echo "True"
    continue
  fi
  within 4 0 3
  if (($v)); then
    echo "True"
    continue
  fi
  within 4 2 1
  if (($v)); then
    echo "True"
  else
    echo "False"
  fi
done <$1
