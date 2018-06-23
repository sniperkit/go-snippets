if [ a -eq b]; then c; fi
if [a -ne b]; then c; fi
if [a -gt b]; then c; fi
if [a -ge b]; then c; fi
if [a -lt b]; then c; fi
if [a -le b]; then c; fi
if ((a < b)); then c; fi
if ((a <= b)); then c; fi
if ((a > b)); then c; fi
if ((a >= b)); then c; fi
if ["a" = "b"]; then c; fi
if [["a" = "b"]]; then c; fi
if ["a" == "b"]; then c; fi
if ["a" == "b*"]; then c; fi
if ["a" != "b"]; then c; fi
if [["a" < "b"]]; then c; fi
if [["a" > "b"]]; then c; fi
if [-Z "a"]; then c; fi
if [-n "a"]; then c; fi
if [a -a b]; then c; fi
if [[a && b]]; then c; fi
if [a -o b]; then c; fi
if [[a || b]]; then c; fi
