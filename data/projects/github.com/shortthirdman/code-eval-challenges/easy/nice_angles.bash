while read line || [ -n "$line" ]; do
  awk 'BEGIN {x='$line'; printf "%d.", int(x); x=(x-int(x))*60; printf "%02d'\''", int(x);  x=(x-int(x))*60; printf "%02d\"\n", int(x)}'
done <$1
