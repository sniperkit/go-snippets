morse="ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90"
while read line || [ -n "$line" ]; do
  m=1
  for ((i=0; i<${#line}; i++)); do
    if [ "${line:$i:1}" = "." ]; then
      ((m*=2))
    elif [ "${line:$i:1}" = "-" ]; then
      ((m=$m*2+1))
    elif [ $m -eq 1 ];then
      printf " "
    else
      if [ $m -lt 64 ]; then
        ((m-=2))
        printf "${morse:m:1}"
      fi
      m=1
    fi
  done
  if [ $m -lt 64 ]; then
    ((m-=2))
    printf "${morse:m:1}"
  fi
  echo
done <$1
