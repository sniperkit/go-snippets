category=("This program is for humans" "Still in Mama's arms" "Preschool Maniac" "Elementary school" "Middle school" "High school" "College" "Working for the man" "The Golden Years")
age=(0 3 5 12 15 19 23 66 101)

while read line || [ -n "$line" ]; do
    m=0
    while [ $m -lt 9 ] && [ $line -ge ${age[$m]} ]; do
        ((m++))
    done
    echo "${category[$m%9]}"
done <$1
