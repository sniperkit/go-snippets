function joins { local d=$1; shift; echo -n "$1"; shift; printf "%s" "${@/#/$d}"; }

primes=(2 3 5 7 11 13)

while read line || [ -n "$line" ]; do
    r=()
    for i in ${primes[*]}; do
        m=$(((1 << $i) - 1))
        if [ $m -ge $line ]; then break; fi
        r+=($m)
    done
    joins ', ' ${r[*]}
    echo
done <$1
