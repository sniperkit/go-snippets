sed -e "s/,,/,x,/g; s/,$/,x/" <$1 | while read line || [ -n "$line" ]; do
  a=( ${line//;/ } )
  if [ ${#a[*]} -lt 2 ]; then
    echo
  else
    b=( ${a[1]//,/ } )
    for ((i=0; i<${#b[*]}; i+=2)); do
      s=${b[$i+1]//0/a}
      s=${s//1/b}
      a=${a//${b[$i]}/$s}
    done
    a=${a//x/}
    a=${a//a/0}
    a=${a//b/1}
    echo "$a"
  fi
done
