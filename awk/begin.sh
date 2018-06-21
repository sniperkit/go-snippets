#statements inside BEGIN are execueted before processing any input text

# FS ---> Filed Seperator
# OFS ---> Output Filed Seperator
echo 'foo:123:bar:789' | awk 'BEGIN{FS=OFS=":"} { print $1, $NF}'
echo 'foo:123:bar:789' | awk -F: -v OFS=":" '{print $1, $NF}'

# changeing a filed will re-build contents of $0
echo ' a    ate b    ' | awk '{$2 = "foo"; print $0}' | cat -A

#$1=$1 is an idiomatic way to re-build when there is nothing else to change

echo 'foo:123:bar:789' | awk -F: -v OFS='-' '{print $0}'
echo 'foo:123:bar:789' | awk -F: -v OFS='-' '{$1=$1; print $0}'
echo 'Sample123string54with908numbers' | awk -F '[0-9]+' '{$1=$1; print $0}'
