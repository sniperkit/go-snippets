chr() {
  printf -v v '%03o' $1
  printf -v v \\$v
}

ord() {
  printf -v v '%d' "'$1"
}

ord 'A'
o=$v
while read a || [ -n "$a" ]; do
  ((a--))
  chr $(($a%26+$o))
  b=
  c=
  d=$v
  if [ $a -ge 26 ]; then
    a=$(($a/26-1))
    chr $(($a%26+$o))
    c=$v
    if [ $a -ge 26 ]; then
      a=$(($a/26-1))
      chr $(($a%26+$o))
      b=$v
    fi
  fi
  echo $b$c$d
done <$1
