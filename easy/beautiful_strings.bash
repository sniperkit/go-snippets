ord() {
  printf -v v '%d' "'$1"
}

tr -cd "[:alpha:]\n" <$1|tr "[:upper:]" "[:lower:]"|while read line || [ -n "$line" ]; do
  declare -a ch
  for ((i=0; i<${#line}; i++)); do
    ord ${line:$i:1}
    ((ch[$v]++))
  done
  IFS=$'\n' so=($(sort -nr <<<"${ch[*]}"))
  n=26
  r=0
  for i in ${so[*]}; do
    ((r+=$n*$i))
    ((n--))
  done
  echo "$r"
  unset ch
done
