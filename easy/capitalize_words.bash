while read line || [ -n "$line" ]; do
    f=1
    for ((i=0; i<${#line}; i++)); do
        a="${line:$i:1}"
        if [ -z "${a//[a-zA-Z]}" ] && [ $f -eq 1 ]; then
            printf ${a^^}
            f=0
        elif [ -z "${a//[ ]}" ]; then
            printf "$a"
            f=1
        else
            printf "$a"
            f=0
        fi
    done
    echo
done < $1
