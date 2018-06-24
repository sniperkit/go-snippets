while read -r line || [ -n "$line" ]; do
  a=( ${line//[^[:alpha:]]/ } )
  echo "${a[*],,}"
done <$1
