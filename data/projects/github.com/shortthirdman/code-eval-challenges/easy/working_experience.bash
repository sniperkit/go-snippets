tr ";-" " " <$1 | while read line || [ -n "$line" ]; do
  a=( $line )
  c=0
  declare -a work
  while [ $c -lt ${#a[*]} ]; do
    case ${a[$c]} in
    Jan) m=1;;
    Feb) m=2;;
    Mar) m=3;;
    Apr) m=4;;
    May) m=5;;
    Jun) m=6;;
    Jul) m=7;;
    Aug) m=8;;
    Sep) m=9;;
    Oct) m=10;;
    Nov) m=11;;
    Dec) m=12;;
    esac
    t0=$(( (${a[$c+1]}-1990)*12+$m ))
    ((c+=2))
    case ${a[$c]} in
    Jan) m=1;;
    Feb) m=2;;
    Mar) m=3;;
    Apr) m=4;;
    May) m=5;;
    Jun) m=6;;
    Jul) m=7;;
    Aug) m=8;;
    Sep) m=9;;
    Oct) m=10;;
    Nov) m=11;;
    Dec) m=12;;
    esac
    t1=$(( (${a[$c+1]}-1990)*12+$m ))
    for ((i=$t0; i<=$t1; i++)); do
      work[$i]=true
    done
    ((c+=2))
  done
  echo "$(( ${#work[*]}/12 ))"
  unset work
done
