sed -e "s/^.*;//; s/,/ /g" $1|while read line || [ -n "$line" ]; do
  a=( $line )
  echo "${a[*]}"|tr " " "\n"|sort|uniq -d
done
