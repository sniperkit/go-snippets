while read line || [ -n "$line" ]; do
  if [ ${#line} -gt 55 ]; then
    for ((i=39; i>=0; i--)); do
      if [ "${line:$i:1}" = " " ]; then
        line=${line:0:$i}
        break
      fi
      if [ ${#line} -gt 55 ]; then
        line=${line:0:40}
      fi
    done
    echo "$line... <Read More>"
  else
    echo "$line"
  fi
done <$1
