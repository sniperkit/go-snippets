declare -a a
DONE=false
until $DONE; do
    read line || DONE=true
    s=( $line )
    if [ ${#a[*]} -gt 0 ]; then
        s[0]=$((${s[0]}+${a[0]}))
    fi
    if [ ${#a[*]} -gt 1 ]; then
        s[${#s[*]}-1]=$((${s[${#s[*]}-1]}+${a[${#a[*]}-1]}))
    fi
    for ((i=1; i<${#s[*]}-1; i++)); do
        if [ ${a[$i-1]} -gt ${a[$i]} ]; then
            ((s[$i]+=${a[$i-1]}))
        else
            ((s[$i]+=${a[$i]}))
        fi
    done
    a=( ${s[*]} )
done <$1
m=${a[0]}
for ((i=1; i<${#a[*]}; i++)); do
    if [ ${a[$i]} -gt $m ]; then
        m=${a[$i]}
    fi
done
echo "$m"
