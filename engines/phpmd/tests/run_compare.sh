actual=$(sh run_tests.sh)
expected=$(cat expected.txt)

if [ "$actual" == "$expected" ]
then
  echo "OK"
else
  echo -e "ERROR\n Actual:"$actual"\nExpected:"$expected
  exit 1
fi
