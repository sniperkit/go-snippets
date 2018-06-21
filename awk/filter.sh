# Each block of statements within {} can be prefiexed by an optional condition so that statements will execute only if condition evaluates to true.
# Condition specified without corresponding statements will lead to printing contents of $0 if condition evalutates to true


# if firest filed exactly matches the string 'apple'
awk '$1=="apple" {print $2}' fruits.txt
awk '{
		 if ($1 == "apple") {
			 print $2
		 }
	 }' fruits.txt


# print first filed if second field > 35
# NR > 1 to avoid the header line.
# NR built-in variable contains record number

awk 'NR>1 && $2>35 {print $1}' fruits.txt
awk '{
		if (NR == 1 || $2 < 35) {
			print $0
		}
	}' fruits.txt

# print header and lines with qty < 35

awk 'NR==1 || $2 < 35' fruits.txt
