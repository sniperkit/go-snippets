awk '{print $0}' fruits.txt
awk '{print $1}' fruits.txt
awk '{print $2}' fruits.txt
awk '{print NF}' fruits.txt
awk '{print $NF}' fruits.txt

printf ' a 		ate b\tc      \n'
printf ' a 		ate b\tc      \n' | awk '{printf $1}'
printf ' a 		ate b\tc      \n' | awk '{printf NF}'


printf ' a 		ate b\tc      \n' | awk -F ' ' '{printf $1}'
printf ' a 		ate b\tc      \n' | awk -F ' ' '{printf NF}'



echo 'apple' | awk -v FS= '{print $1}'
echo 'apple' | awk -v FS= '{print $2}'
echo 'apple' | awk -v FS= '{print $NF}'

