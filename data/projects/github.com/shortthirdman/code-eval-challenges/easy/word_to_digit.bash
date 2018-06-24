while read line || [ -n "$line" ]; do
  a=( ${line//;/ } )
  for i in ${a[*]}; do
    case $i in
    zero ) printf 0;;
    one  ) printf 1;;
    two  ) printf 2;;
    three) printf 3;;
    four ) printf 4;;
    five ) printf 5;;
    six  ) printf 6;;
    seven) printf 7;;
    eight) printf 8;;
    nine ) printf 9;;
    esac
  done
  echo
done <$1
