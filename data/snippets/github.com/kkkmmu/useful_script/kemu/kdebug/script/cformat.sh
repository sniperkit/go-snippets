#!/bin/sh -f

file=$1
echo $file

#Delete CPP comment start with //
sed -i '/^[\t]*\/\//d' $file

# Delete in line comment //
sed -i 's/\/\/[^"]*//' $file

# Delete C Single line comment
sed -i 's/\/\*.*\*\///' $file

# Delete C multi-line comment 
sed -i '/^[[:space:]]*\/\*/,/.*\*\//d' $file

# Remove all empty line
#sed -i '/^[[:space:]]*$/d' $file

# Remove multiple empty line to one empty line
sed -i '/^$/{N;/^\n$/D}' $file
