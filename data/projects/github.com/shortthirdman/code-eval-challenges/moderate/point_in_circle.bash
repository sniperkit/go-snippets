tr -cd "[:digit:]. \n\-" <$1 | while read line || [ -n "$line" ]; do
  a=( $line )
  awk 'BEGIN {x='${a[0]}'-('${a[3]}'); y='${a[1]}'-('${a[4]}'); if (x*x+y*y <= '${a[2]}'*'${a[2]}') {print "true"} else {print "false"}}'
done
