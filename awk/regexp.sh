echo 'Sample123string54with908numbers' | awk -F '[0-9]+' '{print $2}'

echo '{foo} bar=baz' | awk -F '[{}= ]+' '{print $1}'
echo '{foo} bar=baz' | awk -F '[{}= ]+' '{print $2}'
echo '{foo} bar=baz' | awk -F '[{}= ]+' '{print $3}'

# the REGEXP is specified within // and by default acts upon $0
# All lines containing the string 'are'
# same as: grep 'are' poem.txt
awk '/are/' poem.txt

# negating REGEXP, same as: grep -v 'are' poem.txt
awk '!/are/' poem.txt

# same as : grep 'are' poem.txt | grep -v 'so'
awk '/area/ && !/so/' poem.txt
awk '$0 !~ "are"' poem.txt

# lines starting with 'a' or 'b'
awk '/^[ab]/' fruits.txt
awk '$0 ~ "^[ab]"' fruits.txt

# print last field of all lines containing 'are'

awk '/are/{print $NF}' poem.txt

awk '/\/foo\/a\//' paths.txt
awk '$0 ~ "/foo/a/"' paths.txt

# if first field contains 'a'
awk  '$1 ~ /a/' fruits.txt

#if first field contains 'a' and qty > 20
awk '$1 ~ /a/ && $2 > 20' fruits.txt

#if first field does NOT contains 'a'
awk '$1 !~ /a/' fruits.txt
