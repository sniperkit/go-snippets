while read line || [ -n "$line" ]; do
  n=0
  l=${#line}
  line=${line/>>-->/x>}
  while [ $l -gt ${#line} ]; do
    ((n++))
    l=${#line}
    line=${line/>>-->/x>}
  done
  l=${#line}
  line=${line/<--<</x<}
  while [ $l -gt ${#line} ]; do
    ((n++))
    l=${#line}
    line=${line/<--<</x<}
  done
  echo "$n"
done <$1
