sed -r "s|([[:digit:]])([[:alpha:]]+)([[:digit:]])|\3\2\1|g" $1 | while read line || [ -n "$line" ]; do
  echo "$line"
done
