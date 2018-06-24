chr() {
  printf -v v '%03o' $1
  printf -v w '%03o' $2
  printf -v v \\$v\\$w
}

ord() {
  printf -v v '%d' "'$1"
}

while read line || [ -n "$line" ]; do
  ord ${line:0:1}
  l=$(($v-97))
  ord ${line:1:1}
  r=$(($v-49))
  declare -a m
  if [ $l -gt 1 ] && [ $r -gt 0 ]; then
    chr $((97+$l-2)) $((49+$r-1))
    m+=($v)
  fi
  if [ $l -gt 1 ] && [ $r -lt 7 ]; then
    chr $((97+$l-2)) $((49+$r+1))
    m+=($v)
  fi
  if [ $l -gt 0 ] && [ $r -gt 1 ]; then
    chr $((97+$l-1)) $((49+$r-2))
    m+=($v)
  fi
  if [ $l -gt 0 ] && [ $r -lt 6 ]; then
    chr $((97+$l-1)) $((49+$r+2))
    m+=($v)
  fi
  if [ $l -lt 7 ] && [ $r -gt 1 ]; then
    chr $((97+$l+1)) $((49+$r-2))
    m+=($v)
  fi
  if [ $l -lt 7 ] && [ $r -lt 6 ]; then
    chr $((97+$l+1)) $((49+$r+2))
    m+=($v)
  fi
  if [ $l -lt 6 ] && [ $r -gt 0 ]; then
    chr $((97+$l+2)) $((49+$r-1))
    m+=($v)
  fi
  if [ $l -lt 6 ] && [ $r -lt 7 ]; then
    chr $((97+$l+2)) $((49+$r+1))
    m+=($v)
  fi
  echo "${m[*]}"
  unset m
done <$1
